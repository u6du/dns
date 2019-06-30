//usr/bin/env go run "$0" "$@"; exit

package main

import (
	"context"
	"net"
	"sync"
	"time"
)

// TODO 先用UDP DNS找
// 找不到
// 		用DOT找DOT找的时候容忍超时
// 		用 t.cn 测试下是不是支持 ipv6

func LookupTXT(host, nameserver string) string {

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
	if err == nil {
		for i := range li {
			txt := li[i]
			if len(txt) > 0 {
				return txt
			}
		}
	}

	return ""

}

var Ipv4 = []string{
	"180.76.76.76:53",      // 百度
	"223.5.5.5:53",         // 阿里云
	"119.29.29.29:53",      // dnspod
	"1.0.0.1:53",           // cloudflare
	"77.88.8.8:53",         // yandex
	"9.9.9.9:53",           // IBM
	"8.8.8.8:53",           // Google
	"101.7.8.9:53",         // cernet 的官方 DNS
	"208.67.222.222:443",   // OPENDNS
	"114.114.114.114:53",   // 114 dns
	"176.103.130.130:5353", // AdGuard DNS
}

var Ipv6 = []string{
	"[2a00:5a60::ad1:0ff]:5353", // AdGuard DNS
	"[2620:0:ccc::2]:443",       // OPENDNS
	"[2400:da00::6666]:53",      // 百度
	"[2001:da8:202:10::36]:53",  // 北京邮电大学 IPv6 DNS 服务器
	"[2001:da8::666]:53",        // 清华大学 TUNA 协会 IPv6 DNS 服务器
	"[2606:4700:4700::1111]:53", // cloudflare
	"[2a02:6b8::feed:0ff]:53",   //yandex
	"[2001:4860:4860::8888]:53", // Google
	"[2620:fe::fe]:53",          // IBM
}

func ResolveTxt(host string, nameserver []string, verify func(string) bool) string {
	if len(nameserver) == 0 {
		return ""
	}

	ch := make(chan string)

	for i := range nameserver {
		go func(server string) {
			txt := LookupTXT(host, server)
			println("nameserver ", server, " : ", txt)

			ch <- txt
		}(nameserver[i])
	}

	total := 1

	for {
		txt := <-ch
		// TODO 验证签名和时间
		//if len(txt) > 0 || total>=len(nameserver) {
		if verify(txt) || total >= len(nameserver) {
			return txt
		}
		total++
	}
}

func main() {
	//host:="baidu.com"
	//host := "ip4.6du.host"
	host := "g.cn"

	wg := sync.WaitGroup{}

	var v6txt, v4txt string
	run := func(nameserver []string, out *string) {
		wg.Add(1)
		go func() {
			*out = ResolveTxt(host, nameserver, func(txt string) bool {
				return len(txt) > 0
			})
			wg.Done()
		}()
	}

	run(Ipv6, &v6txt)
	run(Ipv4, &v4txt)

	wg.Wait()
	println("v4", v4txt)
	println("v6", v6txt)
	time.Sleep(3 * time.Second)
}
