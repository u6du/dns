//usr/bin/env go run "$0" "$@"; exit

package dns

import (
	"testing"

	"github.com/u6du/config"
)

// TODO 先用UDP DNS找
// 找不到
// 		用DOT找DOT找的时候容忍超时
// 		用 t.cn 测试下是不是支持 ipv6

var HostBootDefault = "6du.host"

func TestNet(t *testing.T) {
	hostPath := "dns/host/boot/"

	v4host := config.File.OneLine(hostPath+"4", "ip4."+HostBootDefault)
	v4txt := ResolveTxtV4(v4host, func(s string) bool {
		t.Log("ip4  ", s)
		return true
	})

	if v4txt != nil {
		t.Log("v4txt ", *v4txt)
	}
	v6host := config.File.OneLine(hostPath+"6", "ip6."+HostBootDefault)
	v6txt := ResolveTxtV6(v6host, func(s string) bool {
		t.Log("ipv4  ", s)
		return true
	})
	if v6txt != nil {
		t.Log("v6txt ", *v6txt)
	}
	t.Log("ipv6 ", TryIpv6())
	t.Log("ipv4 ", TryIpv4())

	//host:="baidu.com"
	//host := "ip4.6du.host"
	/*
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

	*/
}
