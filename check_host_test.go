package spf

import (
	"net"
	"testing"
)

func TestCheckHost(t *testing.T) {
	type expected struct {
		ip     string
		domain string
		res    Result
	}
	tests := []expected{
		{
			"74.6.143.26",
			"yahoo.com",
			ResultPass,
		},
		{
			"35.191.0.1",
			"gmail.com",
			ResultPass,
		},
		{
			"13.110.224.0",
			"colorado.edu",
			ResultPass,
		},
		{
			"0.0.0.0",
			"colorado.edu",
			ResultSoftFail,
		},
	}
	for _, expected := range tests {
		res, err := CheckHost(net.ParseIP(expected.ip), expected.domain)
		if res != expected.res {
			t.Fatalf("Got checkHost(%s, %s)=%s, expected %s (err: %s)", expected.ip, expected.domain, resultToStr(res), resultToStr(expected.res), err)
		}
	}
}
