package main

import "tls-test/proxy"

func main() {
	// ssl server
	//ssl.StartServer()

	// ssl client
	// curl -k -s --cert-type PEM --cert combined.pem https://localhost:9443
	//fmt.Println(ssl.ClientGet())

	// proxy server
	proxy.StartServer()

}
