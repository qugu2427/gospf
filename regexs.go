package spf

import "regexp"

var (
	RgxSpf            *regexp.Regexp = regexp.MustCompile(`^v=spf1`)
	RgxAll            *regexp.Regexp = regexp.MustCompile(`^[\+\-~\?]?all$`)
	RgxIp4            *regexp.Regexp = regexp.MustCompile(`^[\+\-~\?]?ip4:[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)
	RgxIp4Prefixed    *regexp.Regexp = regexp.MustCompile(`^[\+\-~\?]?ip4:[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\/[0-9]{1,2}$`)
	RgxIp6            *regexp.Regexp = regexp.MustCompile(`^[\+\-~\?]?ip6:[0-9a-fA-F:]+$`)
	RgxIp6Prefixed    *regexp.Regexp = regexp.MustCompile(`^[\+\-~\?]?ip6:[0-9a-fA-F:]+\/[0-9]{1,3}$`)
	RgxA              *regexp.Regexp = regexp.MustCompile(`^[\+\-~\?]?a$`)
	RgxAPrefix        *regexp.Regexp = regexp.MustCompile(`^[\+\-~\?]?a\/\d{1,3}$`)
	RgxADomain        *regexp.Regexp = regexp.MustCompile(`^[\+\-~\?]?a:[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9](?:\.[a-zA-Z]{2,})+$`)
	RgxADomainPrefix  *regexp.Regexp = regexp.MustCompile(`^[\+\-~\?]?a:[a-zA-Z0-9-_\.]+\/\d{1,3}`)
	RgxMx             *regexp.Regexp = regexp.MustCompile(`^[\+\-~\?]?mx$`)
	RgxMxPrefix       *regexp.Regexp = regexp.MustCompile(`^[\+\-~\?]?mx\/\d{1,3}$`)
	RgxMxDomain       *regexp.Regexp = regexp.MustCompile(`^[\+\-~\?]?mx:[a-zA-Z0-9-_\.]+`)
	RgxMxDomainPrefix *regexp.Regexp = regexp.MustCompile(`^[\+\-~\?]?mx:[a-zA-Z0-9-_\.]+\/\d{1,3}$`)
	RgxPtr            *regexp.Regexp = regexp.MustCompile(`^[\+\-~\?]?ptr$`)
	RgxPtrDomain      *regexp.Regexp = regexp.MustCompile(`^[\+\-~\?]?ptr:[a-zA-Z0-9-_\.]+$`)
	RgxExists         *regexp.Regexp = regexp.MustCompile(`^[\+\-~\?]?exists:[a-zA-Z0-9-_\.]+$`)
	RgxInclude        *regexp.Regexp = regexp.MustCompile(`^[\+\-~\?]?include:[a-zA-Z0-9-_\.]+$`)
	RgxRedirect       *regexp.Regexp = regexp.MustCompile(`^redirect=[a-zA-Z0-9-_\.]+$`)
	RgxExp            *regexp.Regexp = regexp.MustCompile(`^exp=[a-zA-Z0-9-_\.]+$`)
)
