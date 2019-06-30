//usr/bin/env go run "$0" "$@"; exit

package main

import (
	"context"
	"net"
	"sync"
	"time"
)


func LookupTXT(host,nameserver string) string{
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
	li, err := resolve.LookupTXT(context.Background(), host)
	if err == nil{
		for i := range li {
			txt := li[i]
			if len(txt)>0{
				return txt
			}
		}
	}

	return ""

}

/*
百度公共DNS, "指引"未来
  180.76.76.76
  2400:da00::6666

  阿里云DNS
  223.5.5.5

  dnspod
  119.29.29.29

  cernet 的官方 DNS
  101.7.8.9

  下一代互联网北京研究中心
  240C::6666

  CNNIC IPv6 DNS 服务器
  2001:dc7:1000::1

  清华大学 TUNA 协会 IPv6 DNS 服务器
  2001:da8::666

 */

var Ipv4 = []string {
	"180.76.76.76:53", // 百度
	"223.5.5.5:53", // 阿里云
	"119.29.29.29:53", // dnspod
	"1.1.1.1:53", // cloudflare
	"77.88.8.8:53", // yandex
	"9.9.9.9:53", // IBM
	"8.8.8.8:53", // Google
	"101.7.8.9:53", // cernet 的官方 DNS
}


var Ipv6 = []string{
	"[2001:da8::666]:53", // 清华大学 TUNA 协会 IPv6 DNS 服务器
	"[2606:4700:4700::1111]:53", // cloudflare
	"[2a02:6b8::feed:0ff]:53", //yandex
	"[2001:4860:4860::8888]:53", // Google
	"[2620:fe::fe]:53", // IBM
}


func ResolveTxt(host string, nameserver []string) string{
	if len(nameserver) == 0{
		return ""
	}

	ch := make(chan string)

	for i:=range nameserver{
		go func(server string) {
			txt := LookupTXT(host,server)
			println("nameserver ", server, " : ", txt,"\n")
			ch <- txt
		}(nameserver[i])
	}

	total := 1

	for {
		txt:= <-ch
		//if len(txt) > 0 || total>=len(nameserver) {
		if total>=len(nameserver) {
			return txt
		}
		total++
	}
}


func main() {
	host:="ip4.6du.host"

	wg := sync.WaitGroup{}

	run := func(nameserver []string) string{
		wg.Add(1)
		ch := make(chan string)
		go func() {
			txt:=ResolveTxt(host, nameserver)
			ch<-txt
			wg.Done()
		}()
		txt:= <-ch
		return txt
	}

	v6txt := run(Ipv6)
	v4txt := run(Ipv4)

	wg.Wait()
	println("v4", v4txt)
	println("v6", v6txt)
	time.Sleep(3*time.Second)
}
