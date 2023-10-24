package main

import "net"

func fetchSpfRecords(domain string) (spfRecords []string, err error) {
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		return
	} else {
		for _, txtRecord := range txtRecords {
			if RgxSpf.MatchString(txtRecord) {
				spfRecords = append(spfRecords, txtRecord)
			}
		}
	}
	return
}

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
