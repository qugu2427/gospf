# Go SPF
A sender policy framework module for go. Based off rules at http://www.open-spf.org/SPF_Record_Syntax/. 100% test coverage.

# Example usage
```go
import (
    "net"
    "spf"
)

senderIp := net.ParseIP("0.0.0.0")
senderDomain := "google.com"

// CheckHost is the only exported function in the module:
result, err := spf.CheckHost(senderIp, senderDomain)

switch result {
    case spf.ResultPass:
        fmt.Println("Passed spf check! Accept mail from sender.")
    case spf.ResultFail:
        fmt.Println("Failed spf check! Reject mail from sender.")
    case spf.ResultSoftFail:
        fmt.Println("Soft failed spf check! Reject mail from sender or mark as spam.")
    case spf.ResultNeutral:
        fmt.Println("Neutral spf check result! Reject mail from sender or mark as spam.")
    case spf.ResultNone:
        fmt.Println("This means there was no spf record for the domain. Reject mail from this sender.")
    case ResultPermError:
        fmt.Println("This probably means the spf record for domain was invalid. Reject mail from this sender.")
        fmt.Println(err)
    case ResultTempError:
        fmt.Println("This is a legacy thing. Just treat this as a perm error if you don't know.")
        fmt.Println(err)
    default:
        fmt.Println(err)
}
```

# Things to keep in mind
- The security of SPF relies on the sender domain having a correctly formatted spf record.
- Spf is NOT the only measure needed verify an email.
