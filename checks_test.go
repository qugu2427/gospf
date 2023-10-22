package main

import (
	"net"
	"testing"
)

func TestCheckExists(t *testing.T) {
	domains := map[string]bool{
		"google.com":                    true,
		"www.yahoo.com":                 true,
		"canvas.colorado.edu":           true,
		"fake domain":                   false,
		"fake.thisfakedomainisfake.com": false,
		"D&(WQ*DGW)QGDQD":               false,
		"":                              false,
	}
	for domain, expected := range domains {
		got := checkExists(domain)
		if got != expected {
			t.Fatalf("Got domain=%s exists=%t, expected exists=%t", domain, got, expected)
		}
	}
}

func TestCheckPtr(t *testing.T) {
	type set struct {
		ip     net.IP
		domain string
		hit    bool
	}
	tests := []set{
		{
			net.ParseIP("74.6.143.26"),
			"yahoo.com",
			true,
		},
		{
			net.ParseIP("74.6.143.26"),
			"google.com",
			false,
		},
	}
	for _, exp := range tests {
		hit := checkPtr(exp.ip, exp.domain)
		if hit != exp.hit {
			t.Fatalf("Expected ptr checkPtr(%s, %s) to result in hit=%t, but got hit=%t", exp.ip, exp.domain, exp.hit, hit)
		}
	}

}

// TODO check ip
func TestCheckIp(t *testing.T) {
	type set struct {
		ip    net.IP
		ipStr string
		hit   bool
	}
	tests := []set{
		{
			net.ParseIP("192.168.1.10"),
			"192.168.1.0/24",
			true,
		},
		{
			net.ParseIP("10.0.0.5"),
			"10.0.0.0/8",
			true,
		},
		{
			net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334"),
			"2001:0db8::/32",
			true,
		},
		{
			net.ParseIP("2606:4700:3037::6815:6438"),
			"2606:4700:3037::/48",
			true,
		},
		{
			net.ParseIP("192.168.1.10"),
			"10.0.0.0/8",
			false,
		},
		{
			net.ParseIP("2606:4700:3037::6815:6438"),
			"1.2.3.4",
			false,
		},
	}
	for _, exp := range tests {
		hit := checkIp(exp.ip, exp.ipStr)
		if exp.hit != hit {
			t.Fatalf("Got checkIp(%s, %s)=%t, wanted %t", exp.ip, exp.ipStr, exp.hit, hit)
		}
	}
}

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
			"irrelevant",
			"v=spf1",
			false,
			ResultPass,
		},

		// all
		{
			net.ParseIP("1.2.3.4"),
			"irrelevant",
			"all",
			true,
			ResultPass,
		},
		{
			net.ParseIP("1.2.3.4"),
			"irrelevant",
			"all:",
			false,
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
			false,
			ResultPass,
		},
		{
			net.ParseIP("1.2.3.4"),
			"irrelevant",
			"exists",
			false,
			ResultPass,
		},

		// invalid word
		{
			net.ParseIP("1.2.3.4"),
			"irrelevant",
			"badword",
			false,
			ResultPass,
		},
	}
	for _, test := range tests {
		hit, res := checkWord(
			test.ip,
			test.domain,
			test.word,
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
