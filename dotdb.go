package dns

import (
	"time"

	"github.com/u6du/config"
	"github.com/u6du/ex"
)

func DotTxt(name string, verify func(string) bool) *string {
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

		start := time.Now()

		txt := DotLookupTxt(name, nameserver, 2)

		cost := uint(time.Since(start).Nanoseconds() / 1000000)

		var verified bool

		if txt == nil {
			cost += 20000
			verified = false
		} else {
			verified = verify(*txt)
			if !verified {
				cost += 10000
			}
		}
		costIdLi = append(costIdLi, [2]uint{cost, id})

		if verified {
			c.Close()
			return txt
		}
	}
	return nil
}
