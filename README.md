# Go Sender Policy Framework (rfc7208) Module
A golang sender policy framework verifier module. Designed to be fully compliant with [rfc7208](https://www.rfc-editor.org/rfc/rfc7208).

## Example Usage
```
go get -u github.com/qugu2427/gospf
```
```
package main

import (
	"fmt"
	"net"

	spf "github.com/qugu2427/gospf"
)

func main() {
	senderIp := net.ParseIP("0.0.0.0")
	senderDomain := "example.com"
    senderSender := "bob@example.com"

	result, err := spf.CheckHost(senderIp, senderDomain, senderSender)

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
** WORK IN PROGRESS **

## Known RFC Violation
Below is a list of ways which this module is known to violate [rfc7208](https://www.rfc-editor.org/rfc/rfc7208). Non of these issues are major.

#### 1. The h (HELO/EHLO) macro
This issue only applies when the `h` macro is used AND the `helo` domain does not match the `mailfrom` domain. (this is rare)

[section 7.2 of rfc7208](https://www.rfc-editor.org/rfc/rfc7208#section-7.2) allows the `h` macro to substitute in the `helo/ehlo` domain. Unfortuently, due to an oversight in the rfc, the `helo` domain is not in scope of `CheckHost()` when called on `mailfrom`. To "fix" this issue, we treat the `h` macro the same as a `d` macro. This "fix" will assumes that the helo and mailfrom domains match.

#### 2. Does not support `exp=` modifier
A secondary function is planned to get the explanation if necessary.