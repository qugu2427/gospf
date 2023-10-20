package main

import "net"

func fetchSpfRecords(domain string) (records []string, err error) {
	// pass
}

func CheckHost(ip net.IP, domain string) (res Result, err error) {
	records, err := fetchSpfRecords(domain)
	if err != nil {
		return ResultPermError, err
	}
	for _, record := range records {

	}
}
