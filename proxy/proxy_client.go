package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	PROXY_ADDR1 = "http://foo:bar@localhost:8888" // 代理服务器地址和端口
	TARGET_URL  = "http://www.baidu.com"          // 目标URL
)

func main() {
	// 创建一个http客户端
	client := &http.Client{}

	// 创建代理URL
	proxyURL, err := url.Parse(PROXY_ADDR1)
	if err != nil {
		fmt.Println("Error parsing proxy URL:", err)
		return
	}

	// 设置代理
	client.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	// 创建新的请求
	req, err := http.NewRequest("GET", TARGET_URL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// 处理返回的响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// 输出返回的内容
	fmt.Println(string(body))
}
