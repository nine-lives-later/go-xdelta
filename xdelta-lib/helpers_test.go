package lib

import (
	"testing"
)

func TestStringRoundtrip(t *testing.T) {
	s := "This is the test message!!"

	p := FromString(s)
	s2 := ToString(p, false)

	if s != s2 {
		t.Error("Strings do not match")
	}
}
