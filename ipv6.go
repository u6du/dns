package dns

import "github.com/u6du/config"

var V6 = config.FileLi(
	"dns/ipv6",
	[]string{
		"[2a00:5a60::ad1:0ff]:5353", // AdGuard DNS
		"[2620:0:ccc::2]:443",       // OPENDNS
		"[2400:da00::6666]:53",      // 百度
		"[2001:da8:202:10::36]:53",  // 北京邮电大学 IPv6 DNS 服务器
		"[2001:da8::666]:53",        // 清华大学 TUNA 协会 IPv6 DNS 服务器
		"[2606:4700:4700::1111]:53", // cloudflare
		"[2a02:6b8::feed:0ff]:53",   //yandex
		"[2001:4860:4860::8888]:53", // Google
		"[2620:fe::fe]:53",          // IBM
	})

func TryIpv6() bool {
	return try(V6)
}

func ResolveTxtV6(host string, verify func(string) bool) *string {
	return ResolveTxt(V6, host, verify)
}
