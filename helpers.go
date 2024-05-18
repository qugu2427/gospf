package spf

import (
	"net"
	"strings"
)

// gets spf txt record for a given domain
//
// returns "", nil if no spf record
//
// returns "", <err> if error
func fetchSpfRecord(domain string) (spfRecord string, err error) {
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		return "", err
	}
	for _, txtRecord := range txtRecords {
		if strings.HasPrefix(txtRecord, "v=spf1 ") {
			return txtRecord, nil
		}
	}
	return
}

// Determines the qualifier and returns word without qualifier
func extractQualifier(word string) (res Result, trimmedWord string) {
	if len(word) > 1 {
		switch word[0] {
		case '+':
			return ResultPass, word[1:]
		case '-':
			return ResultFail, word[1:]
		case '~':
			return ResultSoftFail, word[1:]
		case '?':
			return ResultNeutral, word[1:]
		}
	}
	return ResultPass, word
}
