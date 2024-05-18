package spf

import (
	"fmt"
	"net"
	"strings"
)

// Check if a domain resolves. It could resolve to anything.
func checkExists(domain string) (hit bool) {
	_, err := net.LookupIP(domain)
	if err != nil {
		hit = false
	} else {
		hit = true
	}
	return
}

// Checks reverse dns (checks if ip resolves to domain)
func checkPtr(ip net.IP, domain string) (hit bool) {
	domains, err := net.LookupAddr(ip.String())
	if err == nil {
		for _, d := range domains {
			if strings.HasSuffix(d, domain) || strings.HasSuffix(d, domain+".") {
				return true
			}
		}
	}
	return false
}

/*
 Checks if an net.ip is equal to or in range of an ip string.
 The ip string to check can either be absolute (ex "0.0.0.0") or a range (ex "0.0.0.0/0")
 Both ip4 and ip6 work. If ipstr is invalid false will be returned.
*/
func checkIp(ip net.IP, ipStr string) (hit bool) {
	if strings.Contains(ipStr, "/") {
		_, ipRange, err := net.ParseCIDR(ipStr)
		if err != nil {
			hit = false
		} else if ipRange.Contains(ip) {
			hit = true
		} else {
			hit = false
		}
	} else {
		ip2 := net.ParseIP(ipStr)
		if ip2 == nil {
			return false
		}
		hit = ip.Equal(ip2)
	}
	return
}

// Check if domain has A record of given ip
func checkA(ip net.IP, domain string, prefix int) (hit bool) {
	if prefix == -1 {
		aRecs, err := net.LookupIP(domain)
		if err != nil {
			return false
		}
		for _, rec := range aRecs {
			if rec.Equal(ip) {
				return true
			}
		}
	} else {
		aRecs, err := net.LookupIP(domain)
		if err != nil {
			return false
		}
		for _, rec := range aRecs {
			ipStr := fmt.Sprintf("%s/%d", rec, prefix)
			if checkIp(ip, ipStr) {
				return true
			}
		}
	}
	return false
}

// Performs checkA for each mx record
func checkMx(ip net.IP, domain string, prefix int, lookupsLeft *uint8) (hit bool) {
	mxRecs, err := net.LookupMX(domain)
	if err != nil {
		return false
	}
	for _, rec := range mxRecs {
		mxDomain := strings.TrimSuffix(rec.Host, ".")
		if *lookupsLeft <= 0 {
			return false
		}
		if checkA(ip, mxDomain, prefix) {
			return true
		}
		*lookupsLeft--
	}
	return false
}
