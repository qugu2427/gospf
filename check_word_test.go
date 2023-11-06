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
			"irrelevant.com",
			"ip4:74.6.231.20",
			ResultPass,
		},
		{ // negative specific ip4
			"1.2.3.4",
			"irrelevant.com",
			"ip4:74.6.231.20",
			ResultNone,
		},
		{ // positive ranged ip4
			"74.6.128.10",
			"irrelevant.com",
			"ip4:74.6.231.20/16",
			ResultPass,
		},
		{ // negative ranged ip4
			"74.3.128.10",
			"irrelevant.com",
			"ip4:74.6.231.20/16",
			ResultNone,
		},
		{ // bs ip4
			"999.999.999.999",
			"irrelevant.com",
			"ip4:74.6.231.20/16",
			ResultNone,
		},
		{ // bs ip4 word
			"74.6.128.10",
			"irrelevant.com",
			"ip4:999.999.999.999/99",
			ResultNone,
		},

		// ip6
		{ // positive specific ip6
			"2001:4998:24:120d::1:1",
			"irrelevant.com",
			"ip6:2001:4998:24:120d::1:1",
			ResultPass,
		},
		{ // negative specific ip6
			"2001:4998:25:120d::1:1",
			"irrelevant.com",
			"ip6:2001:4998:24:120d::1:1",
			ResultNone,
		},
		{ // positive ranged ip6
			"2001:0db8:85a3:0000:0000:8a2e:0370:7399",
			"irrelevant.com",
			"ip6:2001:0db8:85a3:0000:0000:8a2e:0370:7334/64",
			ResultPass,
		},
		{ // negative ranged ip6
			"2001:0db8:86a3:0000:0000:8a2e:0370:7334",
			"irrelevant.com",
			"ip6:2001:0db8:85a3:0000:0000:8a2e:0370:7334/64",
			ResultNone,
		},
		{ // bs ip6
			"fffffffffff:::::::::",
			"irrelevant.com",
			"ip6:2001:0db8:85a3:0000:0000:8a2e:0370:7334/64",
			ResultNone,
		},
		{ // bs ip6 word
			"2001:4998:24:120d::1:1",
			"irrelevant.com",
			"ip6:fffffffffff:::::::::/64",
			ResultNone,
		},

		// a (depend on yahoo.com)
		{ // positive implied a
			"74.6.231.20",
			"yahoo.com",
			"a",
			ResultPass,
		},
		{ // negative implied a
			"74.6.231.30",
			"yahoo.com",
			"a",
			ResultNone,
		},
		{ // positive ranged a
			"74.6.231.30",
			"yahoo.com",
			"a/24",
			ResultPass,
		},
		{ // negative ranged a
			"74.6.331.30",
			"yahoo.com",
			"a/24",
			ResultNone,
		},
		{ // positive domain a
			"74.6.231.20",
			"irrelevant.com",
			"a:yahoo.com",
			ResultPass,
		},
		{ // negative domain a
			"74.6.231.30",
			"irrelevant.com",
			"a:yahoo.com",
			ResultNone,
		},
		{ // positive domain ranged a
			"74.6.0.0",
			"irrelevant.com",
			"a:yahoo.com/16",
			ResultPass,
		},
		{ // negative domain ranged a
			"74.0.0.0",
			"irrelevant.com",
			"a:yahoo.com/16",
			ResultNone,
		},

		// ptr (these tests depend on yahoo.com)
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

		// TODO mx

		// exists (these tests depend on google.com existing)
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
		res, err := checkWord(
			net.ParseIP(expected.ip),
			expected.domain,
			expected.word,
			[]string{},
		)
		if res != expected.res {
			t.Fatalf(
				"checkWord(%s, %s, %s).res=%s, expected %s (err=%s)",
				expected.ip,
				expected.domain,
				expected.word,
				resultToStr(res),
				resultToStr(expected.res),
				err,
			)
		}
	}

}
