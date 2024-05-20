# Go Sender Policy Framework (rfc7208) Module
A golang sender policy framework verifier module. Designed to be compliant with [rfc7208](https://www.rfc-editor.org/rfc/rfc7208).

## Example Usage
```
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
    senderSender := "bob@example.com"

	result, err := spf.CheckHost(senderIp, senderDomain, senderSender)

	switch result {
	case spf.ResultPass:
		fmt.Println("accept")
	case spf.ResultFail:
		fmt.Println("reject")
	case spf.ResultSoftFail:
		fmt.Println("reject or mark as spam")
	case spf.ResultNeutral:
		// The sender's spf record made no judgement
		fmt.Println("reject or perform other non spf checks")
	case spf.ResultNone:
		// The sender has no spf record
		fmt.Println("reject or perform other non spf checks")
	case spf.ResultPermError:
		// The sender has an invalid spf record
		fmt.Println("reject and log error")
		fmt.Println(err)
	case spf.ResultTempError:
		// Network error while evluating record
		fmt.Println("reject and log error")
		fmt.Println(err)
	default:
		panic(err)
	}
}
```

# Support
- Mechanisms
	* `all` ✅
	* `include` ✅
	* `a` ✅
	* `mx` ✅
	* `ptr` ✅
	* `ip4` ✅
	* `ip6` ✅
	* `exists` ✅
- Modifiers
	* `redirect` ✅
	* `exp` ❌ (no plan to support this since no one uses it)
	* Unknown modifiers will be ignored.
- Macros
	* `s` ✅
	* `l` ✅
	* `o` ✅
	* `d` ✅
	* `i` ✅
	* `p` ✅
	* `v` ✅
	* `h` ⚠️ (will be treated the same as `d`)
		- [section 7.2 of rfc7208](https://www.rfc-editor.org/rfc/rfc7208#section-7.2) allows the `h` macro to substitute in the `helo/ehlo` domain. Unfortuently, due to an oversight in the rfc, the `helo` domain is not always in scope of `CheckHost()` when called on `mailfrom`. To "fix" this issue, we treat the `h` macro the same as a `d` macro. This "fix" will assumes that the helo and mailfrom domains match. This is an issue only in rare cases when the `h` macro is used AND the `helo` domain does not match the `mailfrom` domain.
		
