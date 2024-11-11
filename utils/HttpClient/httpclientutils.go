package HttpClient

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

func GetHttpClientHandle() (*http.Client, error) {
	// 初始化 CookieJar
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Printf("Failed to create cookie jar: %v", err)
		return nil, err
	}
	// 使用自定义 http.Client
	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	return client, nil
}

func NewRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// 设置模拟浏览器的请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	return req, nil
}
