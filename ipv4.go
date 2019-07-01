package dns

import "github.com/u6du/config"

var V4 = config.FileLi(
	"dns/ipv4",
	[]string{
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
	})

func TryIpv4() bool {
	return try(V4)
}

func ResolveTxtV4(host string, verify func(string) bool) *string {
	return ResolveTxt(V4, host, verify)
}
