package spf

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

var hostnameRgx *regexp.Regexp = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-.]*\.\D{2,4}$`)

// Checks that a sender (MAIL FROM or HELO) and ip has permission to send mail from a domain
// (see rfc7208 section 4.1)
//
// Parameters:
//
//	ip: the net.IP of the sender (either ip6 or ip4)
//	domain: the claimed domain of the sender (ex 'colorado.edu' if mail is from 'bob@colorado.edu )
//  sender: the "MAIL FROM" or "HELO" identity, this will likely match the domain
// Returns:
//
//	res: the Result enum (see README for all possible results)
//	err: an error object, only relevant if res = ResultPermError or ResultTempError
func CheckHost(ip net.IP, domain, sender string) (res Result, err error) {
	// invalid domain should return none
	if hostnameRgx.MatchString(domain) {
		return ResultNone, nil
	}

	senderSplit := strings.Split(sender, "@")
	if len(senderSplit) != 2 || len(senderSplit[1]) < 1 {
		return ResultPermError, fmt.Errorf("invalid sender '%s', should be in form 'local@domain'", sender)
	}

	// "If the <sender> has no local-part, substitute the string 'postmaster' for the local-part"
	if senderSplit[0] == "" {
		senderSplit[0] = "postmaster"
	}

	s := session{
		LookupLimit,
		ip,
		domain,
		sender,
		senderSplit[0],
		senderSplit[1],
		"",
	}
	return s.checkHost()
}
