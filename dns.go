package dns

import "github.com/u6du/config"

type Dns struct {
	Nameserver []string
}

func (d *Dns) Txt(host string, verify func(string) bool) *string {
	return ResolveTxt(d.Nameserver, host, verify)
}

func (d *Dns) TxtTest(host string) bool {
	txt := d.Txt(host, func(txt string) bool {
		return true
	})
	return txt != nil
}

var V4 = Dns{config.File.Li(
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
	})}

var V6 = Dns{config.File.Li(
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
	})}
