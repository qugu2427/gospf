package spf

type Result uint8

const (
	ResultPass Result = iota
	ResultNeutral
	ResultSoftFail
	ResultFail
	ResultPermError
	ResultTempError
	ResultNone
)

func (res Result) String() (str string) {
	switch res {
	case ResultPass:
		str = "Pass"
	case ResultFail:
		str = "Fail"
	case ResultSoftFail:
		str = "SoftFail"
	case ResultNeutral:
		str = "Neutral"
	case ResultNone:
		str = "None"
	case ResultPermError:
		str = "PermError"
	case ResultTempError:
		str = "TempError"
	default:
		panic("unknown result in String()")
	}
	return
}
