//usr/bin/env go run "$0" "$@"; exit

package dns

import (
	"testing"
)

func TestNet(t *testing.T) {
	DotTxt("6.ip.6du.host", func(s string) bool {
		t.Log(s)
		return false
	})
}
