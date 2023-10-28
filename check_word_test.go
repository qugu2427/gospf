package spf

import (
	"net"
	"testing"
)

func TestCheckWord(t *testing.T) {
	type expected struct {
		ip     string
		domain string
		word   string
		res    Result
	}
	tests := []expected{
		// v=spf1
		{
			"1.2.3.4",
			"irrelevant.com",
			"v=spf1",
			ResultNone,
		},

		// all
		{
			"1.2.3.4",
			"irrelevant.com",
			"all",
			ResultPass,
		},

		// ip4
		{ // positive specific ip4
			"74.6.231.20",
			"yahoo.com",
			"ip4:74.6.231.20",
			ResultPass,
		},
		{ // negative specific ip4
			"1.2.3.4",
			"yahoo.com",
			"ip4:74.6.231.20",
			ResultNone,
		},
		{ // positive ranged ip4
			"74.6.128.10",
			"yahoo.com",
			"ip4:74.6.231.20/16",
			ResultPass,
		},
		{ // negative ranged ip4
			"74.3.128.10",
			"yahoo.com",
			"ip4:74.6.231.20/16",
			ResultNone,
		},

		// ip6
		{ // positive specific ip6
			"2001:4998:24:120d::1:1",
			"yahoo.com",
			"ip6:2001:4998:24:120d::1:1",
			ResultPass,
		},
		{ // negative specific ip6
			"2001:4998:25:120d::1:1",
			"yahoo.com",
			"ip6:2001:4998:24:120d::1:1",
			ResultNone,
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
			ResultPass,
		},
		{ // negative implied ptr
			"0.0.0.0",
			"yahoo.com",
			"ptr",
			ResultNone,
		},
		{ // positive specified ptr
			"74.6.231.20",
			"irrelevant.com",
			"ptr:yahoo.com",
			ResultPass,
		},
		{ // negative specified ptr
			"0.0.0.0",
			"irrelevant",
			"ptr:yahoo.com",
			ResultNone,
		},

		// exists
		{ // positive exists
			"1.2.3.4",
			"irrelevant",
			"exists:google.com",
			ResultPass,
		},
		{ // negative exists (non existent domain)
			"1.2.3.4",
			"irrelevant",
			"exists:spf.thisdomainistotallyfake.gov",
			ResultNone,
		},
		{ // negative/err exists (no domain specified)
			"1.2.3.4",
			"irrelevant",
			"exists",
			ResultPermError,
		},

		// invalid word
		{
			"1.2.3.4",
			"irrelevant.com",
			"invalidword",
			ResultPermError,
		},
	}
	for _, expected := range tests {
		res, _ := checkWord(
			net.ParseIP(expected.ip),
			expected.domain,
			expected.word,
			[]string{},
		)
		if res != expected.res {
			t.Fatalf(
				"checkWord(%s, %s, %s).res=%s, expected %s",
				expected.ip,
				expected.domain,
				expected.word,
				resultToStr(res),
				resultToStr(expected.res),
			)
		}
	}

}
