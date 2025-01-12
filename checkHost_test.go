package spf

import (
	"net"
	"testing"

	"github.com/foxcpp/go-mockdns"
)

var dnsMock map[string]mockdns.Zone = map[string]mockdns.Zone{
	"basic.ip4.": {
		TXT: []string{"v=spf1 ip4:1.2.3.4 ?ip4:1.2.3.5 ~ip4:1.2.3.6 -all"},
	},

	"bad.spf.": {
		TXT: []string{"v=spf3 ip4:1.2.3.4 ?ip4:1.2.3.5 ~ip4:1.2.3.6 -all"},
	},

	"multi.include.": {
		TXT: []string{"v=spf1 ~include:sub.multi.include -all"},
	},
	"sub.multi.include.": {
		TXT: []string{"v=spf1 include:pass.all -all"},
	},

	"pass.all.": {
		TXT: []string{"v=spf1 +all"},
	},

	"bad.include.": {
		TXT: []string{"v=spf1 include:no.where ~all"},
	},

	"bad.redirect.": {
		TXT: []string{"v=spf1 redirect=no.where"},
	},
	"good.redirect.": {
		TXT: []string{"v=spf1 redirect=pass.all"},
	},
	"ignore.redirect.": {
		TXT: []string{"v=spf1 redirect=pass.all -all"},
	},
}

func TestCheckHost(t *testing.T) {
	srv, _ := mockdns.NewServer(dnsMock, false)
	defer srv.Close()
	srv.PatchNet(net.DefaultResolver)
	defer mockdns.UnpatchNet(net.DefaultResolver)

	type test struct {
		ip             net.IP
		domain         string
		sender         string
		expectedResult Result
		expectErr      bool
	}

	tests := []test{

		// Basic ip4 qualifier tests
		{
			net.ParseIP("1.2.3.4"),
			"basic.ip4",
			"test@basic.ip4",
			ResultPass,
			false,
		},
		{
			net.ParseIP("1.2.3.5"),
			"basic.ip4",
			"test@basic.ip4",
			ResultNeutral,
			false,
		},
		{
			net.ParseIP("1.2.3.6"),
			"basic.ip4",
			"test@basic.ip4",
			ResultSoftFail,
			false,
		},
		{
			net.ParseIP("0.0.0.0"),
			"basic.ip4",
			"test@basic.ip4",
			ResultFail,
			false,
		},

		// DNS errors (such as non-existent domain) give TempError
		{
			net.ParseIP("0.0.0.0"),
			"fakedomain.xyzzz",
			"test@fakedomain.xyzzz",
			ResultTempError,
			true,
		},

		// A domain with no spf1 record should give None
		{
			net.ParseIP("0.0.0.0"),
			"bad.spf",
			"test@bad.spf",
			ResultNone,
			false,
		},

		// Two layer include
		{
			net.ParseIP("0.0.0.0"),
			"multi.include",
			"test@multi.include",
			ResultSoftFail,
			false,
		},

		// An include to nowhere should just ignore it
		{
			net.ParseIP("0.0.0.0"),
			"bad.include",
			"test@bad.include",
			ResultSoftFail,
			false,
		},

		// A redirect to nowhere should be a permerror
		{
			net.ParseIP("0.0.0.0"),
			"bad.redirect",
			"test@bad.redirect",
			ResultPermError,
			true,
		},

		// A valid redirect test
		{
			net.ParseIP("0.0.0.0"),
			"good.redirect",
			"test@good.redirect",
			ResultPass,
			false,
		},

		// Redirect should always be evaluated last
		{
			net.ParseIP("0.0.0.0"),
			"ignore.redirect",
			"test@ignore.redirect",
			ResultFail,
			false,
		},
	}

	for _, test := range tests {
		gotResult, err := CheckHost(test.ip, test.domain, test.sender)
		if err != nil && !test.expectErr {
			t.Fatalf("got unexpected err from CheckHost(%s, %s, %s), %s", test.ip, test.domain, test.sender, err)
		} else if err == nil && gotResult != test.expectedResult {
			t.Fatalf("got unexpected result %s from CheckHost(%s, %s, %s) expected %s", gotResult, test.ip, test.domain, test.sender, test.expectedResult)
		}
	}

}
