# Go SPF
A super simple sender policy framework module for go. Based off rules at http://www.open-spf.org/SPF_Record_Syntax/.

# Example usage
```shell
go get -u github.com/qugu2427/gospf
```
```go
package main

import (
	"fmt"
	"net"

	spf "github.com/qugu2427/gospf"
)

func main() {
	senderIp := net.ParseIP("0.0.0.0")
	senderDomain := "example.com"

	result, err := spf.CheckHost(senderIp, senderDomain)

	switch result {
	case spf.ResultPass:
		fmt.Println("accept")
	case spf.ResultFail:
		fmt.Println("reject")
	case spf.ResultSoftFail:
		fmt.Println("accept but mark")
	case spf.ResultNeutral:
		fmt.Println("accept")
	case spf.ResultNone:
		fmt.Println("accept")
	case spf.ResultPermError:
		fmt.Println("unspecified")
		fmt.Println(err)
	case spf.ResultTempError:
		fmt.Println("accept or reject")
		fmt.Println(err)
	default:
		panic(err)
	}
}
```

# Things to keep in mind
- The security of SPF relies on the sender domain having a correctly formatted spf record.
- Spf is NOT the only measure needed verify an email.
- There are other spf modules in go, but most of the ones I've seen are overly complicated, messy, incomplete, and untested.
- The 'exp' modifier is ignored & not handled by this module because no one uses it and it does not add anything to the authentication. SPF records with exp will work fine, you will just have to go get the explanation domain yourself.
