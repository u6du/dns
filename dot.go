package main

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/rs/zerolog/log"
)

func NewResolver(addr string) *net.Resolver {
	var dialer net.Dialer
	tlsConfig := &tls.Config{
		ServerName:         addr,
		ClientSessionCache: tls.NewLRUClientSessionCache(32),
		InsecureSkipVerify: false,
	}

	return &net.Resolver{
		PreferGo: true,
		Dial: func(context context.Context, _, address string) (net.Conn, error) {
			conn, err := dialer.DialContext(context, "tcp", addr+":853")
			if err != nil {
				return nil, err
			}

			_ = conn.(*net.TCPConn).SetKeepAlive(true)
			_ = conn.(*net.TCPConn).SetKeepAlivePeriod(10 * time.Minute)
			return tls.Client(conn, tlsConfig), nil
		},
	}

}

func Txt(name, nameserver string, retry int) string {

	var resolve *net.Resolver

	if len(nameserver) > 0 {
		resolve = NewResolver(nameserver)
	} else {
		resolve = net.DefaultResolver
	}

	n := 1

	for {
		li, err := resolve.LookupTXT(context.Background(), name)
		if err != nil {
			log.Warn().Err(err).Msg(nameserver)

			if n >= retry {
				break
			}

			n += 1
		}

		for i := range li {
			return li[i]
		}
	}
	return ""
}

/*
func init() {
	p2p.DefaultResolver = newResolver("1.0.0.1")
}
*/

/*
func DialNew(nameserver string) func(context.Context, string, string) (p2p.Conn, error) {
	return func(ctx context.Context, network, address string) (p2p.Conn, error) {
		d := p2p.Dialer{}
		return d.DialContext(ctx, "udp", nameserver+":53")
	}
}

var NAMESERVER = []string{
	"180.76.76.76",
	"119.29.29.29",
	"223.5.5.5",
	"1.1.1.1",
	"208.67.220.220",
	"8.8.8.8",
	"9.9.9.9",
	"114.114.114.114",
	"223.6.6.6",
}

func main() {
	for _, nameserver := range NAMESERVER {


	}
}
*/
