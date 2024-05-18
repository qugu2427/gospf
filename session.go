package spf

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Section 4.6.4 limits the amount of dns query requests to 10
var LookupLimit uint8 = 10

// An object to track each spf lookup prodecure
type session struct {
	lookupsLeft  uint8
	ip           net.IP
	domain       string
	sender       string
	senderLocal  string
	senderDomain string
	redirect     string
}

func (s *session) checkHost() (res Result, err error) {
	// TODO check dns name

	// Get spf records / TempError if err
	record, err := fetchSpfRecord(s.domain)
	if err != nil {
		return ResultTempError, err
	}

	// Return ResultNone if no spf record
	if record == "" {
		return ResultNone, nil
	}

	words := strings.Split(record, " ")
	var hit bool
	l := len(words)
	for i := 1; i < l; i++ { // skip first word ("v=spf1")

		if s.lookupsLeft <= 0 {
			return ResultPermError, fmt.Errorf("lookup limit exceeded")
		}

		hit, res, err = s.checkWord(words[i])
		if err != nil {
			return ResultPermError, err
		}
		if hit {
			return
		}
	}

	// If a redirect has been specified
	// go to it once all other words have been evaled
	if s.redirect != "" {
		s.domain = s.redirect
		s.redirect = ""
		res, err = s.checkHost()
	} else {
		res = ResultNeutral
	}

	return
}

func (s *session) checkWord(word string) (hit bool, result Result, err error) {
	result, word = extractQualifier(word)
	if strings.Contains(word, "%") {
		word, err = applyMacros(word, s)
	}
	l := len(word)
	if word == "all" {
		hit = true
	} else if l > 4 && (strings.HasPrefix(word, "ip4:") || strings.HasPrefix(word, "ip6:")) {
		ipStr := word[4:]
		hit = checkIp(s.ip, ipStr)
	} else if word == "a" {
		hit = checkA(s.ip, s.domain, -1)
		s.lookupsLeft--
	} else if l > 2 && strings.HasPrefix(word, "a/") {
		var prefix int
		prefix, err = strconv.Atoi(word[2:])
		if err == nil {
			hit = checkA(s.ip, s.domain, prefix)
			s.lookupsLeft--
		}
	} else if l > 2 && strings.HasPrefix(word, "a:") {
		if strings.Contains(word, "/") {
			var prefix int
			slashIndex := strings.Index(word, "/")
			prefix, err = strconv.Atoi(word[slashIndex+1:])
			if err == nil {
				hit = checkA(s.ip, word[2:slashIndex], prefix)
			}
		} else {
			domain := word[2:]
			hit = checkA(s.ip, domain, -1)
		}
		s.lookupsLeft--
	} else if word == "mx" {
		hit = checkMx(s.ip, s.domain, -1, &s.lookupsLeft)
	} else if l > 3 && strings.HasPrefix(word, "mx/") {
		var prefix int
		prefix, err = strconv.Atoi(word[3:])
		if err == nil {
			hit = checkMx(s.ip, s.domain, prefix, &s.lookupsLeft)
		}
	} else if l > 3 && strings.HasPrefix(word, "mx:") {
		if strings.Contains(word, "/") {
			var prefix int
			slashIndex := strings.Index(word, "/")
			prefix, err = strconv.Atoi(word[slashIndex+1:])
			if err == nil {
				hit = checkMx(s.ip, word[2:slashIndex], prefix, &s.lookupsLeft)
			}
		} else {
			domain := word[2:]
			hit = checkMx(s.ip, domain, -1, &s.lookupsLeft)
		}
	} else if word == "ptr" {
		hit = checkPtr(s.ip, s.domain)
		s.lookupsLeft--
	} else if l > 4 && strings.HasPrefix(word, "ptr:") {
		domain := word[4:]
		hit = checkPtr(s.ip, domain)
		s.lookupsLeft--
	} else if l > 7 && strings.HasPrefix(word, "exists:") {
		domain := word[7:]
		hit = checkExists(domain)
		s.lookupsLeft--
	} else if l > 8 && strings.HasPrefix(word, "include:") {
		domain := word[8:]
		sCopy := session{s.lookupsLeft, s.ip, domain, s.sender, s.senderLocal, s.senderDomain, ""}
		copyRes, copyErr := sCopy.checkHost()
		if copyErr == nil && copyRes == ResultPass {
			hit = true
		}
	} else if l > 9 && strings.HasPrefix(word, "redirect=") {
		s.redirect = word[9:]
		s.lookupsLeft--
	} else if !strings.Contains(word, "=") || l < 2 {
		err = fmt.Errorf("unknown word '%s'", word)
	}
	return
}
