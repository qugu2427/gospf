package main

import (
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

// Checks reverse dns (if ip has A record of domain)
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

func checkWord(ip net.IP, domain, word string) (hit bool, res Result) {
	res = getQualifierResult(word)
	switch {
	case RgxSpf.MatchString(word):
		hit = false
	case RgxAll.MatchString(word):
		hit = true
	case RgxIp4.MatchString(word):
		ipStr := word[strings.Index(word, ":")+1:]
		hit = checkIp(ip, ipStr)
	case RgxIp4Prefixed.MatchString(word):
		ipStr := word[strings.Index(word, ":")+1:]
		hit = checkIp(ip, ipStr)
	case RgxIp6.MatchString(word):
		ipStr := word[strings.Index(word, ":")+1:]
		hit = checkIp(ip, ipStr)
	case RgxIp6Prefixed.MatchString(word):
		ipStr := word[strings.Index(word, ":")+1:]
		hit = checkIp(ip, ipStr)
	case RgxA.MatchString(word):
		// pass
	case RgxAPrefix.MatchString(word):
		// pass
	case RgxADomain.MatchString(word):
		// pass
	case RgxADomainPrefix.MatchString(word):
		// pass
	case RgxMx.MatchString(word):
		// pass
	case RgxMxPrefix.MatchString(word):
		// pass
	case RgxMxDomainPrefix.MatchString(word):
		// pass
	case RgxPtr.MatchString(word): // ptr
		hit = checkPtr(ip, domain)
	case RgxPtrDomain.MatchString(word): // ptr:<domain>
		domain = word[strings.Index(word, ":")+1:]
		hit = checkPtr(ip, domain)
	case RgxExists.MatchString(word): // exists:<domain>
		domain = word[strings.Index(word, ":")+1:]
		hit = checkExists(domain)
	case RgxInclude.MatchString(word): // include:<comain>
		// pass
	case RgxRedirect.MatchString(word): // redirect=<domain>
		// pass
	case RgxExp.MatchString(word): // exp=<domain>
		// pass
	default:
		hit = false
		res = ResultPermError
	}
	return
}

func CheckHost(ip net.IP, domain string) (res Result, err error) {
	records, err := fetchSpfRecords(domain)
	if err != nil {
		return ResultPermError, err
	}
	for _, record := range records {
		words := strings.Split(record, " ")
		for _, word := range words {
			hit, res := checkWord(ip, domain, word)
			if hit {
				return res, err
			}
		}
	}
	res = ResultNone
	return
}
