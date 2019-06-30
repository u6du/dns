package main

import (
	"fmt"
	"net"

	"github.com/u6du/ex"
)

func main() {
	dnsServer := "8.8.8.8"
	c, err := net.Dial("udp", dnsServer)
	ex.Panic(err)

	fmt.Println("vim-go")
}
