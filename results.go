package spf

type Result = uint8

// Results are in order of greatness
// e.x [ Fail, Fail, Pass, Fail ] = Pass
const (
	ResultPass uint8 = iota
	ResultNeutral
	ResultSoftFail
	ResultFail
	ResultPermError
	ResultTempError
	ResultNone
)

func resultToStr(res Result) (str string) {
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
		panic("unknown result in resultToStr()")
	}
	return
}

func getDominantResult(resArr []Result) (res Result) {
	res = ResultNone
	for _, r := range resArr {
		if r < res {
			res = r
		}
	}
	return
}
