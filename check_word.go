package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func checkWord(ip net.IP, domain, word string, domainsVisited []string) (hit bool, res Result, err error) {
	fmt.Printf("[DEBUG] %#v -> ", word)
	res = getQualifierResult(word)
	switch {
	case RgxSpf.MatchString(word): // v=spf1
		hit = false
	case RgxAll.MatchString(word): // all
		hit = true
	case RgxIp4.MatchString(word): // ip4:<ip>
		ipStr := word[4:]
		hit = checkIp(ip, ipStr)
	case RgxIp4Prefixed.MatchString(word): // ip4:<ip>/<prefix>
		ipStr := word[4:]
		hit = checkIp(ip, ipStr)
	case RgxIp6.MatchString(word): // ip6:<domain>
		ipStr := word[4:]
		hit = checkIp(ip, ipStr)
	case RgxIp6Prefixed.MatchString(word): // ip6:<domain>/<prefix>
		ipStr := word[4:]
		hit = checkIp(ip, ipStr)
	case RgxA.MatchString(word): // a
		hit = checkA(ip, domain, -1)
	case RgxAPrefix.MatchString(word): // a/<prefix>
		prefix, err := strconv.Atoi(word[3:])
		if err != nil {
			hit = false
		} else {
			hit = checkA(ip, domain, prefix)
		}
	case RgxADomain.MatchString(word): // a:<domain>
		domain := word[3:]
		hit = checkA(ip, domain, -1)
	case RgxADomainPrefix.MatchString(word): // a:<domain>/<prefix>
		slashIndex := strings.Index(word, "/")
		prefix, err := strconv.Atoi(word[slashIndex+1:])
		if err != nil {
			hit = false
		} else {
			hit = checkA(ip, word[3:slashIndex], prefix)
		}
	case RgxMx.MatchString(word): // mx
		hit = checkMx(ip, domain, -1)
	case RgxMxPrefix.MatchString(word): // mx/<prefix>
		prefix, err := strconv.Atoi(word[3:])
		if err != nil {
			hit = false
		} else {
			hit = checkMx(ip, domain, prefix)
		}
	case RgxMxDomainPrefix.MatchString(word): // mx:domain/<prefix>
		slashIndex := strings.Index(word, "/")
		prefix, err := strconv.Atoi(word[slashIndex+1:])
		if err != nil {
			hit = false
		} else {
			hit = checkMx(ip, word[3:slashIndex], prefix)
		}
	case RgxPtr.MatchString(word): // ptr
		hit = checkPtr(ip, domain)
	case RgxPtrDomain.MatchString(word): // ptr:<domain>
		domain = word[4:]
		hit = checkPtr(ip, domain)
	case RgxExists.MatchString(word): // exists:<domain>
		domain = word[7:]
		hit = checkExists(domain)
	case RgxInclude.MatchString(word): // include:<domain>
		includeDomain := word[8:]
		res, err = checkHostInner(ip, domain, append(domainsVisited, includeDomain))
	case RgxRedirect.MatchString(word): // redirect=<domain>
		redirectDomain := word[9:]
		fmt.Println(redirectDomain)
		res, err = checkHostInner(ip, domain, append(domainsVisited, redirectDomain))
	case RgxExp.MatchString(word): // exp=<domain>
		hit = false
		// pass
	case strings.Contains(word, "="): // unknown modifier
		hit = false
	default:
		hit = true
		res = ResultPermError
	}
	fmt.Printf("%t %#v\n", hit, res)
	return
}
