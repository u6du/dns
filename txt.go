package dns

import (
	"context"
	"net"
)

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
		if verify(txt) || total >= len(nameserver) {
			return txt
		}
		total++
	}
}
