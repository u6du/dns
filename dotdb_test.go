//usr/bin/env go run "$0" "$@"; exit

package dns

import (
	"testing"
)

func TestDotdb(t *testing.T) {
	DotTxt("4.ip.6du.host", func(s string) bool {
		t.Log(s)
		return false
	})
}
