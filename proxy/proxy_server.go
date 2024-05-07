package proxy

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// curl -x http://foo:bar@127.0.0.1:8888  http://www.baidu.com

const (
	PROXY_ADDR     = "172.16.16.235:8888"
	PROXY_USERNAME = "foo" // 替换为你的代理用户名
	PROXY_PASSWORD = "bar" // 替换为你的代理密码
)

func handleRequestAndRedirect(rw http.ResponseWriter, req *http.Request) {
	// 检查代理授权头
	proxyAuth := req.Header.Get("Proxy-Authorization")
	if proxyAuth == "" {
		rw.Header().Set("Proxy-Authenticate", "Basic realm=\"Proxy\"")
		rw.WriteHeader(http.StatusProxyAuthRequired)
		return
	}

	// 验证凭据
	if !checkProxyAuth(proxyAuth) {
		http.Error(rw, "Proxy Authentication Required", http.StatusProxyAuthRequired)
		return
	}

	log.Printf("start to request...%s\n", req.RequestURI)
	for key, value := range req.Header {
		log.Printf("===request===, name: %s, value: %s\n", key, value[0])
	}
	// 修改请求信息以发往目的服务器
	req.URL.Scheme = "http"
	req.URL.Host = req.Host
	req.RequestURI = ""

	// 转发请求到目的服务器
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadGateway)
		return
	}
	defer response.Body.Close()

	for key, value := range response.Header {
		log.Printf("===response===, name: %s, value: %s\n", key, value[0])
	}
	copyHeader(rw.Header(), response.Header)
	rw.WriteHeader(response.StatusCode)
	io.Copy(rw, response.Body)
}

func checkProxyAuth(proxyAuth string) bool {
	basicPrefix := "Basic "
	if !strings.HasPrefix(proxyAuth, basicPrefix) {
		return false
	}

	encodedCreds := strings.TrimPrefix(proxyAuth, basicPrefix)
	creds, err := base64.StdEncoding.DecodeString(encodedCreds)
	if err != nil {
		return false
	}

	parts := strings.SplitN(string(creds), ":", 2)
	if len(parts) != 2 {
		return false
	}

	username := parts[0]
	password := parts[1]

	// 简单的用户名和密码校验
	return username == PROXY_USERNAME && password == PROXY_PASSWORD
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func StartServer() {
	http.HandleFunc("/", handleRequestAndRedirect)
	fmt.Println("Serving proxy on :8888...")
	if err := http.ListenAndServe(PROXY_ADDR, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
