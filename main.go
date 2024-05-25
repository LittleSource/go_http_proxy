package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// 处理函数，将所有请求转发到目标网址
func handler(w http.ResponseWriter, r *http.Request) {
	// 构建目标 URL
	targetURL := "http://yixiang.qccyou.com" + r.RequestURI

	// 创建一个新的请求
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 复制请求头
	req.Header = make(http.Header)
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// 发送请求并获取响应
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 将目标服务器的响应内容发送给客户端
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func main() {
	// 设置请求处理函数
	http.HandleFunc("/", handler)

	// 启动 HTTP 服务器
	fmt.Println("HTTP server is running on http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}
