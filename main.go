package main

import (
	"fmt"
	"tls-test/ssl"
)

func main() {
	//ssl.StartServer()
	fmt.Println(ssl.ClientGet())
}
