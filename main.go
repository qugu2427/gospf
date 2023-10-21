package main

import "fmt"

func main() {
	records, err := fetchSpfRecords("google.com")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(records)
}
