package main

import (
	"encoding/binary"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/u6du/ex"
	"github.com/u6du/go-rfc1924/base85"
	"golang.org/x/crypto/ed25519"

	"sixdu/config"
	"sixdu/key"
)

const (
	Success uint8 = iota
	ErrVerify
	ErrTimeout
	ErrDecode
	ErrEmpty
)

const TimeOutHour = uint32(8)

func Decode(txt string) ([]byte, error) {
	b, err := base85.DecodeString(txt)
	return b, err
}

func Parse(txt string) ([]byte, uint8) {
	if len(txt) > ed25519.SignatureSize {
		b, err := Decode(txt)
		if err != nil {
			return []byte{}, ErrDecode
		}
		n := ed25519.SignatureSize
		ctx := b[n:]
		sign := b[:n]
		if ed25519.Verify(key.GodPublic, ctx, sign) {

			hour := ctx[0:4]
			ctx := ctx[4:]

			cost := uint32(time.Now().Unix()/3600) - binary.LittleEndian.Uint32(hour)
			if cost >= TimeOutHour {
				return ctx, ErrTimeout
			} else {
				return ctx, Success
			}

		} else {
			return []byte{}, ErrVerify
		}
	}
	return []byte{}, ErrEmpty
}

type Resolve struct {
	timeoutCount uint8
}

func (r *Resolve) Txt(name, nameserver string, retry int) ([]byte, uint8) {
	txt := dns.Txt(name, nameserver, retry)
	b, state := Parse(txt)
	if state == ErrTimeout {
		if r.timeoutCount < 3 {
			r.timeoutCount += 1
		} else {
			return b, Success
		}
	}
	return b, state
}

func BootNode(name string) []byte {
	resolve := Resolve{}

	b, state := resolve.Txt(name, "", 1)
	if state == Success {
		return b
	}

	db := config.Db(
		"dns/dot",

		`CREATE TABLE "dot" (
"id"	INTEGER PRIMARY KEY AUTOINCREMENT,
"host"	TEXT NOT NULL UNIQUE,
"delay"	INTEGER NOT NULL DEFAULT 0);
CREATE INDEX "dot.delay" ON "dot" ("delay" ASC);`,

		"INSERT INTO dot(host) values (?)",

		"dns.rubyfish.cn",
		"dot-jp.blahdns.com",
		"dns.google",
		"security-filter-dns.cleanbrowsing.org",
		"dot.securedns.eu",
		"sdns.233py.com",
		"edns.233py.com",
		"ndns.233py.com",
		"dns.quad9.net",
		"wdns.233py.com",
		"dot-de.blahdns.com",
		"1dot1dot1dot1.cloudflare-dns.com",
		"dns.brahma.world",
	)

	defer db.Close()

	c, err := db.Query("select id,host from dot order by delay asc")

	ex.Panic(err)

	var id uint
	var nameserver string
	var costIdLi [][2]uint

	defer func() {
		for _, costId := range costIdLi {
			_, err := db.Exec("UPDATE dot SET delay=? WHERE id=?", costId[0], costId[1])
			ex.Panic(err)
		}
	}()

	for c.Next() {
		c.Scan(&id, &nameserver)
		log.Debug().Msg(nameserver)
		start := time.Now()

		b, state = resolve.Txt(name, nameserver, 2)

		cost := uint(time.Since(start).Nanoseconds() / 1000000)

		success := state == Success

		if !success {
			cost += 99999
		}
		costIdLi = append(costIdLi, [2]uint{cost, id})

		if success {
			c.Close()
			return b
		}
	}
	return []byte{}
}
