package dns

import (
	"context"
	"net"
)

func LookupTXT(host, nameserver string) *string {

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
			return &li[i]
		}
	}

	return nil

}

func ResolveTxt(host string, nameserver []string, verify func(*string) bool) *string {
	if len(nameserver) == 0 {
		return nil
	}

	ch := make(chan *string)
	defer close(ch)

	for i := range nameserver {
		go func(server string) {
			txt := LookupTXT(host, server)
			if txt != nil {
				println("nameserver ", server, " : ", *txt)
			}

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
	return nil
}
