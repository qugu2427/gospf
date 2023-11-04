package spf

import (
	"fmt"
	"net"
	"strings"
)

const IncludeDepthLimit int = 100

func checkHostInner(ip net.IP, domain string, domainsVisited []string) (res Result, err error) {
	dprint("calling checkHostInner(%s, %s, %#v)", ip, domain, domainsVisited)

	// Make sure we are not stuck in a loop or bottomless spf search
	if hasDuplicateDomain(domainsVisited, domain) {
		res = ResultPermError
		err = fmt.Errorf("spf record loop detected over domain %s, likely circular redirect or include in the domains spf records", domain)
		return
	} else if len(domainsVisited) > IncludeDepthLimit {
		res = ResultPermError
		err = fmt.Errorf("spf record depth limit reached over domain %s, likely too many includes or redirects in the domains spf records", domain)
		return
	}

	// Get spf records / PermError if err
	record, err := fetchSpfRecord(domain)
	if err != nil {
		return ResultPermError, err
	}

	// Parse spf record to array of results
	var results []Result
	words := strings.Split(record, " ")
	for _, word := range words {
		res, err = checkWord(ip, domain, word, domainsVisited)
		results = append(results, res)
	}

	res = getDominantResult(results)
	return
}

// Checks that a sender ip has permission to send mail from a domain
//
// Parameters:
// 	ip: the net.IP of the sender (either ip6 or ip4)
// 	domain: the claimed domain of the sender (ex 'colorado.edu' if mail is from 'bob@colorado.edu )
//
// Returns:
// 	res: the Result enum (see README for all possible results)
// 	err: and error object, only relevant if res = ResultPermError or ResultTempError
func CheckHost(ip net.IP, domain string) (res Result, err error) {
	return checkHostInner(ip, domain, []string{})
}
