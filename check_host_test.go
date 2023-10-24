package main

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
	}
	for _, expected := range tests {
		res, err := CheckHost(net.ParseIP(expected.ip), expected.domain)
		if res != expected.res {
			t.Fatalf("Got checkHost(%s, %s)=%#v, expected %#v (err: %s)", expected.ip, expected.domain, res, expected.res, err)
		}
	}
}
