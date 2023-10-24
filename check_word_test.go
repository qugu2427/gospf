package main

import (
	"net"
	"testing"
)

func TestCheckWord(t *testing.T) {
	type set struct {
		ip     net.IP
		domain string
		word   string
		expHit bool
		expRes Result
	}
	tests := []set{
		// v=spf1
		{
			net.ParseIP("1.2.3.4"),
			"irrelevant.com",
			"v=spf1",
			false,
			ResultPass,
		},

		// all
		{
			net.ParseIP("1.2.3.4"),
			"irrelevant.com",
			"all",
			true,
			ResultPass,
		},

		// ip4
		{
			net.ParseIP("74.6.231.20"),
			"yahoo.com",
			"ip4:74.6.231.20",
			true,
			ResultPass,
		},
		{
			net.ParseIP("74.6.128.10"),
			"yahoo.com",
			"ip4:74.6.231.20/16",
			true,
			ResultPass,
		},
		{
			net.ParseIP("1.2.3.4"),
			"yahoo.com",
			"ip4:74.6.231.20",
			false,
			ResultPass,
		},
		{
			net.ParseIP("74.3.128.10"),
			"yahoo.com",
			"ip4:74.6.231.20/16",
			false,
			ResultPass,
		},

		// ip6
		{
			net.ParseIP("2001:4998:24:120d::1:1"),
			"yahoo.com",
			"ip6:2001:4998:24:120d::1:1",
			true,
			ResultPass,
		},

		// ptr
		{
			net.ParseIP("74.6.231.20"),
			"yahoo.com",
			"ptr",
			true,
			ResultPass,
		},
		{
			net.ParseIP("0.0.0.0"),
			"yahoo.com",
			"ptr",
			false,
			ResultPass,
		},
		{
			net.ParseIP("74.6.231.20"),
			"irrelevant",
			"ptr:yahoo.com",
			true,
			ResultPass,
		},
		{
			net.ParseIP("0.0.0.0"),
			"irrelevant",
			"ptr:yahoo.com",
			false,
			ResultPass,
		},

		// exists
		{
			net.ParseIP("1.2.3.4"),
			"irrelevant",
			"exists:google.com",
			true,
			ResultPass,
		},
		{
			net.ParseIP("1.2.3.4"),
			"irrelevant",
			"exists:bad",
			true,
			ResultPermError,
		},
		{
			net.ParseIP("1.2.3.4"),
			"irrelevant",
			"exists",
			true,
			ResultPermError,
		},

		// invalid word
		{
			net.ParseIP("1.2.3.4"),
			"irrelevant.com",
			"invalidword",
			true,
			ResultPermError,
		},
	}
	for _, test := range tests {
		hit, res, _ := checkWord(
			test.ip,
			test.domain,
			test.word,
			[]string{},
		)
		if hit != test.expHit {
			t.Fatalf(
				"Expected hit=%t on checkWord(%s, %s, %s), got hit=%t",
				test.expHit,
				test.ip,
				test.domain,
				test.word,
				hit,
			)
		} else if hit && res != test.expRes {
			t.Fatalf(
				"Expected res=%#v on checkWord(%s, %s, %s), got res=%#v",
				test.expRes,
				test.ip,
				test.domain,
				test.word,
				res,
			)
		}
	}

}
