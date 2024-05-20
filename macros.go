package spf

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

var (
	macroRgx      *regexp.Regexp = regexp.MustCompile(`%{[^ %{}]*}`)
	macroBodyRgx  *regexp.Regexp = regexp.MustCompile(`^[slodipvh][0-9]{0,3}r?[\.\-+,/_=]{0,7}$`)
	macroDigitRgx *regexp.Regexp = regexp.MustCompile(`[0-9]{1,3}`)
)

func applyMacros(word string, s *session) (macroWord string, err error) {
	word = strings.ReplaceAll(word, "%%", "%")
	word = strings.ReplaceAll(word, "%_", " ")
	word = strings.ReplaceAll(word, "%-", "%20")
	matches := macroRgx.FindAllString(word, -1)
	for _, match := range matches {
		macroBody := match[2 : len(match)-1]
		macroResult, err := parseMacro(macroBody, s)
		if err != nil {
			return "", err
		}
		word = strings.ReplaceAll(word, match, macroResult)
	}
	return word, nil
}

// trims to right sections of a '.' delimenated word
// e.x: one.two.three, 2 --> two.three
func macroTrim(word string, rightTrim int) (tword string) {
	chars := []rune(word)
	l := len(chars) - 1
	for i := l; i >= 0; i-- {
		if chars[i] == '.' {
			rightTrim--
			if rightTrim <= 0 {
				return
			}
		}
		tword = string(chars[i]) + tword
	}
	return
}

// reverses a '.' deliminated word
// e.x: one.two.three -> three.two.one
func macroReverse(word string) (rword string) {
	subWord := ""
	for _, char := range word {
		subWord += string(char)
		if char == '.' {
			rword = "." + subWord[:len(subWord)-1] + rword
			subWord = ""
		}
	}
	rword = subWord + rword

	return
}

func parseMacro(macroBody string, s *session) (word string, err error) {
	macroBody = strings.ToLower(macroBody)
	if !macroBodyRgx.MatchString(macroBody) {
		return "", fmt.Errorf("invalid macro body `%s`", macroBody)
	}

	// Extract right trim digits from macro
	rightTrim := 0
	rightTrimStr := macroDigitRgx.FindString(macroBody)
	if rightTrimStr != "" {
		rightTrim, err = strconv.Atoi(rightTrimStr)
		if err != nil {
			return "", err
		}
		if rightTrim > 127 {
			return "", fmt.Errorf("right trim too large")
		}
		macroBody = macroDigitRgx.ReplaceAllString(macroBody, "")
	}

	// Substitute macro data
	switch macroBody[0] {
	case 's':
		word = s.sender
	case 'l':
		word = s.senderLocal
	case 'o':
		word = s.senderDomain
	case 'd':
		word = s.domain
	case 'i':
		word = s.ip.String()
	case 'p':
		domains, err := net.LookupAddr(s.ip.String())
		if err != nil {
			return "", err
		}
		if len(domains) != 1 {
			return "", fmt.Errorf("p macro evaluation resulted in %d domains", len(domains))
		}
		word = domains[0]
	case 'v':
		if s.ip.To4() != nil {
			word = "in-addr"
		} else {
			word = "ip6"
		}
	case 'h': // KNOWN BUG: Due to oversight in rfc, h is just the domain
		word = s.domain
	default:
		return "", fmt.Errorf("invalid macro letter '%s'", string(macroBody[0]))
	}

	// Determine reverse and delimeters
	delimeters := []string{}
	reverse := false
	l := len(macroBody)
	for i := 1; i < l; i++ {
		if macroBody[i] == 'r' {
			reverse = true
		} else {
			delimeters = append(delimeters, string(macroBody[i]))
		}
	}

	// Replace all specified delims with '.'
	for _, delim := range delimeters {
		word = strings.ReplaceAll(word, delim, ".")
	}

	// Reverse and trim
	if reverse {
		word = macroReverse(word)
	}
	if rightTrim != 0 {
		word = macroTrim(word, rightTrim)
	}

	return
}
