package spf

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
			t.Fatalf("Got checkExists(%s)=%t, expected %t", domain, got, expected)
		}
	}
}

func TestCheckPtr(t *testing.T) {
	type expected struct {
		ip     string
		domain string
		hit    bool
	}
	tests := []expected{
		{
			"74.6.143.26",
			"yahoo.com",
			true,
		},
		{
			"74.6.143.26",
			"google.com",
			false,
		},
	}
	for _, expected := range tests {
		hit := checkPtr(net.ParseIP(expected.ip), expected.domain)
		if hit != expected.hit {
			t.Fatalf("Got checkPtr(%s, %s)=%t, expected %t", expected.ip, expected.domain, hit, expected.hit)
		}
	}

}

func TestCheckIp(t *testing.T) {
	type expected struct {
		ip    string
		ipStr string
		hit   bool
	}
	tests := []expected{
		{
			"192.168.1.10",
			"192.168.1.0/24",
			true,
		},
		{
			"10.0.0.5",
			"10.0.0.0/8",
			true,
		},
		{
			"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			"2001:0db8::/32",
			true,
		},
		{
			"2606:4700:3037::6815:6438",
			"2606:4700:3037::/48",
			true,
		},
		{
			"192.168.1.10",
			"10.0.0.0/8",
			false,
		},
		{
			"2606:4700:3037::6815:6438",
			"1.2.3.4",
			false,
		},
	}
	for _, expected := range tests {
		hit := checkIp(net.ParseIP(expected.ip), expected.ipStr)
		if expected.hit != hit {
			t.Fatalf("Got checkIp(%s, %s)=%t, expected %t", expected.ip, expected.ipStr, hit, expected.hit)
		}
	}
}

func TestCheckA(t *testing.T) {
	type expected struct {
		ip     string
		domain string
		prefix int
		hit    bool
	}
	tests := []expected{
		{
			"74.6.231.20",
			"yahoo.com",
			-1,
			true,
		},
		{
			"2001:4998:24:120d::1:1",
			"yahoo.com",
			-1,
			true,
		},
		{
			"0.0.0.0",
			"yahoo.com",
			-1,
			false,
		},
		{
			"74.6.123.45",
			"yahoo.com",
			16,
			true,
		},
		{
			"2001:4998:124:1507::f000",
			"yahoo.com",
			32,
			true,
		},
		{
			"74.5.123.45",
			"yahoo.com",
			16,
			false,
		},
		{
			"2002:4998:124:1507::f000",
			"yahoo.com",
			32,
			false,
		},
	}
	for _, expected := range tests {
		hit := checkA(net.ParseIP(expected.ip), expected.domain, expected.prefix)
		if hit != expected.hit {
			t.Fatalf("Got checkA(%s, %s %d)=%t, expected %t", expected.ip, expected.domain, expected.prefix, hit, expected.hit)
		}
	}
}

// These dont always pass since mx changes alot
// FIXME BUG
// func TestCheckMx(t *testing.T) {
// 	type expected struct {
// 		ip     string
// 		domain string
// 		prefix int
// 		hit    bool
// 	}
// 	tests := []expected{
// 		{
// 			"67.195.228.110",
// 			"yahoo.com",
// 			-1,
// 			true,
// 		},
// 		{
// 			"0.0.0.0",
// 			"yahoo.com",
// 			-1,
// 			false,
// 		},
// 		{
// 			"67.195.228.0",
// 			"yahoo.com",
// 			24,
// 			true,
// 		},
// 		{
// 			"67.195.227.0",
// 			"yahoo.com",
// 			24,
// 			false,
// 		},
// 	}
// 	for _, expected := range tests {
// 		hit := checkMx(net.IP(expected.ip), expected.domain, expected.prefix)
// 		if hit != expected.hit {
// 			t.Fatalf("Got checkMx(%s, %s %d)=%t, expected %t", expected.ip, expected.domain, expected.prefix, hit, expected.hit)
// 		}
// 	}
// }
