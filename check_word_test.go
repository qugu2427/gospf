package main

import (
	"net"
	"testing"
)

func TestCheckWord(t *testing.T) {
	type expected struct {
		ip     string
		domain string
		word   string
		hit    bool
		res    Result
	}
	tests := []expected{
		// v=spf1
		{
			"1.2.3.4",
			"irrelevant.com",
			"v=spf1",
			false,
			ResultPass,
		},

		// all
		{
			"1.2.3.4",
			"irrelevant.com",
			"all",
			true,
			ResultPass,
		},

		// ip4
		{ // positive specific ip4
			"74.6.231.20",
			"yahoo.com",
			"ip4:74.6.231.20",
			true,
			ResultPass,
		},
		{ // negative specific ip4
			"1.2.3.4",
			"yahoo.com",
			"ip4:74.6.231.20",
			false,
			ResultPass,
		},
		{ // positive ranged ip4
			"74.6.128.10",
			"yahoo.com",
			"ip4:74.6.231.20/16",
			true,
			ResultPass,
		},
		{ // negative ranged ip4
			"74.3.128.10",
			"yahoo.com",
			"ip4:74.6.231.20/16",
			false,
			ResultPass,
		},

		// ip6
		{ // positive specific ip6
			"2001:4998:24:120d::1:1",
			"yahoo.com",
			"ip6:2001:4998:24:120d::1:1",
			true,
			ResultPass,
		},
		{ // negative specific ip6
			"2001:4998:25:120d::1:1",
			"yahoo.com",
			"ip6:2001:4998:24:120d::1:1",
			false,
			ResultPass,
		},
		// TODO positive ranged ip6
		// TODO negative ranged ip6

		// a
		// TODO positive implied a
		// TODO negative implied a
		// TODO positive ranged a
		// TODO negative ranged a
		// TODO positive domain a
		// TODO negative domain a
		// TODO positive domain ranged a
		// TODO negative domain ranged a

		// ptr
		{ // positive implied ptr
			"74.6.231.20",
			"yahoo.com",
			"ptr",
			true,
			ResultPass,
		},
		{ // negative implied ptr
			"0.0.0.0",
			"yahoo.com",
			"ptr",
			false,
			ResultPass,
		},
		{ // positive specified ptr
			"74.6.231.20",
			"irrelevant.com",
			"ptr:yahoo.com",
			true,
			ResultPass,
		},
		{ // negative specified ptr
			"0.0.0.0",
			"irrelevant",
			"ptr:yahoo.com",
			false,
			ResultPass,
		},

		// exists
		{ // positive exists
			"1.2.3.4",
			"irrelevant",
			"exists:google.com",
			true,
			ResultPass,
		},
		{ // negative exists (non existent domain)
			"1.2.3.4",
			"irrelevant",
			"exists:spf.thisdomainistotallyfake.gov",
			false,
			ResultPermError,
		},
		{ // negative exists (invalid domain format)
			"1.2.3.4",
			"irrelevant",
			"exists:bad",
			true,
			ResultPermError,
		},
		{ // negative exists (no domain specified)
			"1.2.3.4",
			"irrelevant",
			"exists",
			true,
			ResultPermError,
		},

		// invalid word
		{
			"1.2.3.4",
			"irrelevant.com",
			"invalidword",
			true,
			ResultPermError,
		},
	}
	for _, expected := range tests {
		hit, res, _ := checkWord(
			net.ParseIP(expected.ip),
			expected.domain,
			expected.word,
			[]string{},
		)
		if hit != expected.hit {
			t.Fatalf(
				"checkWord(%s, %s, %s).hit=%t, expected %t",
				expected.ip,
				expected.domain,
				expected.word,
				hit,
				expected.hit,
			)
		} else if hit && res != expected.res {
			t.Fatalf(
				"checkWord(%s, %s, %s).res=%#v, expected %#v",
				expected.ip,
				expected.domain,
				expected.word,
				res,
				expected.res,
			)
		}
	}

}
