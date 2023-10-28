package spf

import (
	"fmt"
	"net"
)

// Gets spf txt record for a given domain
// only 1 spf record is allowed
func fetchSpfRecord(domain string) (spfRecord string, err error) {
	txtRecords, err := net.LookupTXT(domain)
	hasFoundSpfRecord := false
	if err != nil {
		return
	} else {
		for _, txtRecord := range txtRecords {
			if RgxSpf.MatchString(txtRecord) {
				if hasFoundSpfRecord {
					return "", fmt.Errorf("more than one spf record found for %s", domain)
				} else {
					spfRecord = txtRecord
					hasFoundSpfRecord = true
				}
			}
		}
	}
	if !hasFoundSpfRecord {
		err = fmt.Errorf("no spf record found for %s", domain)
	}
	return
}

// Determines the qualifier of a word and returns it's results
func getQualifierResult(word string) Result {
	switch word[0] {
	case '-':
		return ResultFail
	case '~':
		return ResultSoftFail
	case '?':
		return ResultNeutral
	}
	return ResultPass
}

func hasDuplicateDomain(domainsVisited []string, domain string) bool {
	isFound := false
	for _, d := range domainsVisited {
		if d == domain && isFound {
			return true
		} else if d == domain {
			isFound = true
		}
	}
	return false
}

// This is just a way to debug print
// DEV ONLY
var shouldDebugPrint bool = false

func dprint(msg string, fmtArgs ...interface{}) {
	if shouldDebugPrint {
		str := fmt.Sprintf(msg, fmtArgs...)
		fmt.Println("[DEBUG] " + str)
	}
}
