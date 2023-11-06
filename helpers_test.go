package spf

import "testing"

func TestGetQualifierResult(t *testing.T) {
	words := map[string]Result{
		"-word": ResultFail,
		"~word": ResultSoftFail,
		"?word": ResultNeutral,
		"+word": ResultPass,
		"word":  ResultPass,
	}
	for word, expected := range words {
		got := getQualifierResult(word)
		if got != expected {
			t.Fatalf("Got getQualifierResult(%s)=%s, expected %s", word, resultToStr(got), resultToStr(expected))
		}
	}
}

func TestHasDuplicateDomain(t *testing.T) {
	type expected struct {
		domains      []string
		domain       string
		hasDuplicate bool
	}
	tests := []expected{
		{
			[]string{},
			"x",
			false,
		},
		{
			[]string{"x"},
			"x",
			false,
		},
		{
			[]string{"x", "y", "z", "w", "x"},
			"x",
			true,
		},
		{
			[]string{"hit", "y", "z", "w", "hit", "hit", "hit"},
			"hit",
			true,
		},
	}
	for _, test := range tests {
		hasDuplicate := hasDuplicateDomain(test.domains, test.domain)
		if hasDuplicate != test.hasDuplicate {
			t.Fatalf("Got hasDuplicateDomain(%#v, %s)=%t, expected %t", test.domains, test.domain, hasDuplicate, test.hasDuplicate)
		}
	}
}
