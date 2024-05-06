package ssl

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// curl -k -s --cert-type PEM --cert combined.pem https://localhost:9443
func ClientGet() string {
	// 设置客户端证书
	//certPath := "combined.pem"
	// 从PEM文件中读取客户端证书和私钥
	cert, err := tls.LoadX509KeyPair("ssl/combined.pem", "ssl/combined.pem")
	if err != nil {
		log.Fatalf("failed to load client certificate: %v", err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: true,
			},
		},
	}

	url := "https://127.0.0.1:9443"
	//req, err := http.NewRequest("HEAD", url, nil)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Failed to create request:", err)
		return ""
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send request:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response:", err)
		os.Exit(1)
	}
	return string(body)
}
