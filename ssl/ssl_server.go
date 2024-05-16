package ssl

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
)

func StartServer() {
	// 加载CA证书
	caCert, err := os.ReadFile("ssl/ca.pem")
	if err != nil {
		log.Fatalf("读取CA证书失败: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// 创建TLS配置
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert, // 需要并验证客户端证书
	}
	tlsConfig.BuildNameToCertificate()

	// 创建HTTP服务器并设置TLS配置
	server := &http.Server{
		Addr:      ":9443",
		TLSConfig: tlsConfig,
	}

	// 设定路由
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("recevied request, %s", r.Header.Get("user-agent"))
		w.Header().Set("Set-cookie", "aaa=bbb")
		w.Write([]byte("Hello, this is a secure server!"))
	})

	// 监听并在TLS层服务
	log.Printf("服务器启动，监听 :9443")
	log.Fatal(server.ListenAndServeTLS("ssl/server.pem", "ssl/server.key")) // 使用服务端证书和密钥文件
}
