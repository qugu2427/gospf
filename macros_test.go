package spf

import (
	"net"
	"testing"
)

func TestApplyMacros(t *testing.T) {
	type test struct {
		session   session
		word      string
		expected  string
		expectErr bool
	}

	s := session{
		10,
		net.ParseIP("192.0.2.3"),
		"email.example.com",
		"strong-bad@email.example.com",
		"strong-bad",
		"email.example.com",
		"",
	}

	tests := []test{
		{
			s,
			"%{s}",
			"strong-bad@email.example.com",
			false,
		},
		{
			s,
			"%{o}",
			"email.example.com",
			false,
		},
		{
			s,
			"%{d}",
			"email.example.com",
			false,
		},
		{
			s,
			"%{d4}",
			"email.example.com",
			false,
		},
		{
			s,
			"%{d3}",
			"email.example.com",
			false,
		},
		{
			s,
			"%{d2}",
			"example.com",
			false,
		},
		{
			s,
			"%{d1}",
			"com",
			false,
		},
		{
			s,
			"%{dr}",
			"com.example.email",
			false,
		},
		{
			s,
			"%{d2r}",
			"example.email",
			false,
		},
		{
			s,
			"%{l}",
			"strong-bad",
			false,
		},
		{
			s,
			"%{l-}",
			"strong.bad",
			false,
		},
		{
			s,
			"%{lr}",
			"strong-bad",
			false,
		},
		{
			s,
			"%{lr-}",
			"bad.strong",
			false,
		},
		{
			s,
			"%{l1r-}",
			"strong",
			false,
		},
		{
			s,
			"%{ir}.%{v}._spf.%{d2}",
			"3.2.0.192.in-addr._spf.example.com",
			false,
		},
		{
			s,
			"%{lr-}.lp._spf.%{d2}",
			"bad.strong.lp._spf.example.com",
			false,
		},
		{
			s,
			"%{lr-}.lp.%{ir}.%{v}._spf.%{d2}",
			"bad.strong.lp.3.2.0.192.in-addr._spf.example.com",
			false,
		},
		{
			s,
			"%{ir}.%{v}.%{l1r-}.lp._spf.%{d2}",
			"3.2.0.192.in-addr.strong.lp._spf.example.com",
			false,
		},
		{
			s,
			"%{d2}.trusted-domains.example.net",
			"example.com.trusted-domains.example.net",
			false,
		},
	}

	for _, test := range tests {
		got, err := applyMacros(test.word, &test.session)
		if err != nil && !test.expectErr {
			t.Fatalf("got unexpected err '%s' on applyMacros(%s, %#v)", err, test.word, &test.session)
		} else if got != test.expected {
			t.Fatalf("got '%s' on applyMacros(%s, %#v), expected %s", got, test.word, test.session, test.expected)
		}
	}
}

func TestMacroReverse(t *testing.T) {
	type test struct {
		word     string
		expected string
	}

	tests := []test{
		{
			"",
			"",
		},
		{
			"word.",
			".word",
		},
		{
			"one.two",
			"two.one",
		},
		{
			"one.two.three",
			"three.two.one",
		},
		{
			"one.two.three.",
			".three.two.one",
		},
		{
			"one.two.three..",
			"..three.two.one",
		},
	}

	for _, test := range tests {
		got := macroReverse(test.word)
		if got != test.expected {
			t.Fatalf("got '%s' macroReverse(%s), expected %s", got, test.word, test.expected)
		}
	}
}

func TestMacroTrim(t *testing.T) {
	type test struct {
		word      string
		trimRight int
		expected  string
	}

	tests := []test{
		{
			"",
			1,
			"",
		},
		{
			"one.two.three",
			4,
			"one.two.three",
		},
		{
			"one.two.three",
			3,
			"one.two.three",
		},
		{
			"two.three",
			2,
			"two.three",
		},
		{
			"three",
			1,
			"three",
		},
	}

	for _, test := range tests {
		got := macroTrim(test.word, test.trimRight)
		if got != test.expected {
			t.Fatalf("got '%s' macroTrim(%s, %d), expected %s", got, test.word, test.trimRight, test.expected)
		}
	}
}
