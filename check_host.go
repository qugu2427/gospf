package main

import (
	"fmt"
	"net"
	"strings"
)

func checkHostInner(ip net.IP, domain string, domainsVisited []string) (res Result, err error) {
	fmt.Printf("[DEBUG] calling checkhostinner on %s\n", domain)

	// Make sure we are not stuck in a loop or bottomless spf search
	if hasDuplicateDomain(domainsVisited, domain) {
		res = ResultPermError
		err = fmt.Errorf("spf record loop detected over domain %s, likely circular redirect or include in the domains spf records", domain)
		return
	} else if len(domainsVisited) > 100 {
		res = ResultPermError
		err = fmt.Errorf("spf record depth limit reached over domain %s, likely too many includes or redirects in the domains spf records", domain)
		return
	}

	records, err := fetchSpfRecords(domain)
	fmt.Printf("[DEBUG] %#v\n", records)
	if err != nil {
		return ResultPermError, err
	}
	for _, record := range records {
		words := strings.Split(record, " ")
		var hit bool
		for _, word := range words {
			hit, res, err = checkWord(ip, domain, word, domainsVisited)
			if hit {
				return res, err
			}
		}
	}
	res = ResultNone
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
	fmt.Printf("[DEBUG] calling checkhost on %s\n", domain)
	return checkHostInner(ip, domain, []string{})
}
