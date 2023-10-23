package main

import (
	"net"
	"strconv"
	"strings"
)

func checkWord(ip net.IP, domain, word string) (hit bool, res Result) {
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
	case RgxIp6.MatchString(word):
		ipStr := word[4:]
		hit = checkIp(ip, ipStr)
	case RgxIp6Prefixed.MatchString(word):
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
	case RgxMx.MatchString(word):
		hit = checkMx(ip, domain, -1) // mx
	case RgxMxPrefix.MatchString(word): // mx/<prefix>
		prefix, err := strconv.Atoi(word[3:])
		if err != nil {
			hit = false
		} else {
			hit = checkA(ip, domain, prefix)
		}
	case RgxMxDomainPrefix.MatchString(word): // mx:domain/<prefix>
		slashIndex := strings.Index(word, "/")
		prefix, err := strconv.Atoi(word[slashIndex+1:])
		if err != nil {
			hit = false
		} else {
			hit = checkA(ip, word[3:slashIndex], prefix)
		}
	case RgxPtr.MatchString(word): // ptr
		hit = checkPtr(ip, domain)
	case RgxPtrDomain.MatchString(word): // ptr:<domain>
		domain = word[4:]
		hit = checkPtr(ip, domain)
	case RgxExists.MatchString(word): // exists:<domain>
		domain = word[7:]
		hit = checkExists(domain)
	case RgxInclude.MatchString(word): // include:<comain>
		// pass
	case RgxRedirect.MatchString(word): // redirect=<domain>
		// pass
	case RgxExp.MatchString(word): // exp=<domain>
		// pass
	default:
		hit = true
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
