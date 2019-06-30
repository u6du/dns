//usr/bin/env go run "$0" "$@"; exit

package main

import (
	"context"
	"fmt"
	"net"

	"github.com/u6du/ex"
)


func LookupTXT(nameserver string){
	resolve := &net.Resolver{
		PreferGo: true,
		Dial: func(context context.Context, _, address string) (net.Conn, error) {
			var dialer net.Dialer

			conn, err := dialer.DialContext(context, "udp", nameserver)
			if err != nil {
				return nil, err
			}

			return conn, nil
		},
	}
	li, err := resolve.LookupTXT(context.Background(), "ip4.6du.host")
	ex.Panic(err)
	fmt.Printf("%s", li)
}

func main() {
	LookupTXT("8.8.8.8:53")
	LookupTXT("[2620:fe::fe]:53")
}
