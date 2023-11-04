package spf

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func checkWord(ip net.IP, domain, word string, domainsVisited []string) (res Result, err error) {
	res = ResultNone
	switch {
	case RgxAll.MatchString(word): // all
		res = getQualifierResult(word)
	case RgxIp4.MatchString(word) ||
		RgxIp4Prefixed.MatchString(word) ||
		RgxIp6.MatchString(word) ||
		RgxIp6Prefixed.MatchString(word): // ip (any format)
		ipStr := word[4:]
		if checkIp(ip, ipStr) {
			res = getQualifierResult(word)
		} else {
			res = ResultNone
		}
	case RgxA.MatchString(word): // a
		if checkA(ip, domain, -1) {
			res = getQualifierResult(word)
		}
	case RgxAPrefix.MatchString(word): // a/<prefix>
		var prefix int
		prefix, err = strconv.Atoi(word[3:])
		if err != nil {
			res = ResultPermError
		} else if checkA(ip, domain, prefix) {
			res = getQualifierResult(word)
		}
	case RgxADomain.MatchString(word): // a:<domain>
		domain := word[3:]
		if checkA(ip, domain, -1) {
			res = getQualifierResult(word)
		}
	case RgxADomainPrefix.MatchString(word): // a:<domain>/<prefix>
		var prefix int
		slashIndex := strings.Index(word, "/")
		prefix, err = strconv.Atoi(word[slashIndex+1:])
		if err != nil {
			res = ResultPermError
		} else if checkA(ip, word[3:slashIndex], prefix) {
			res = getQualifierResult(word)
		}
	case RgxMx.MatchString(word): // mx
		if checkMx(ip, domain, -1) {
			res = getQualifierResult(word)
		}
	case RgxMxPrefix.MatchString(word): // mx/<prefix>
		var prefix int
		prefix, err = strconv.Atoi(word[3:])
		if err != nil {
			res = ResultPermError
		} else if checkMx(ip, domain, prefix) {
			res = getQualifierResult(word)
		}
	case RgxMxDomainPrefix.MatchString(word): // mx:domain/<prefix>
		var prefix int
		slashIndex := strings.Index(word, "/")
		prefix, err = strconv.Atoi(word[slashIndex+1:])
		if err != nil {
			res = ResultPermError
		} else if checkMx(ip, word[3:slashIndex], prefix) {
			res = getQualifierResult(word)
		}
	case RgxPtr.MatchString(word): // ptr
		if checkPtr(ip, domain) {
			res = getQualifierResult(word)
		}
	case RgxPtrDomain.MatchString(word): // ptr:<domain>
		domain = word[4:]
		if checkPtr(ip, domain) {
			res = getQualifierResult(word)
		}
	case RgxExists.MatchString(word): // exists:<domain>
		domain = word[8:]
		if checkExists(domain) {
			res = getQualifierResult(word)
		}
	case RgxInclude.MatchString(word): // include:<domain>
		includeDomain := word[8:]
		res, err = checkHostInner(ip, includeDomain, append(domainsVisited, includeDomain))
	case RgxRedirect.MatchString(word): // redirect=<domain>
		redirectDomain := word[9:]
		res, err = checkHostInner(ip, redirectDomain, append(domainsVisited, redirectDomain))
	case RgxExp.MatchString(word): // exp=<domain>
		// TODO
	case !strings.Contains(word, "="): // invalid word (excluding unknown modifiers)
		res = ResultPermError
		err = fmt.Errorf("invalid word '%s'", word)
	}
	dprint("\t %s -> %#v", word, res)
	return
}
