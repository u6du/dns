package dns

import (
	"time"

	"github.com/u6du/zerolog/log"
	"github.com/u6du/ex"

	"github.com/u6du/config"
)

type Resolve struct {
	timeoutCount uint8
}

func (r *Resolve) Txt(name, nameserver string, retry int) ([]byte, uint8) {
	txt := Txt(name, nameserver, retry)
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
