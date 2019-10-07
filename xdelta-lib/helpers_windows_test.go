// +build windows

package lib

import (
	"testing"
)

func TestStringRoundtrip(t *testing.T) {
	s := "This is the test message!!"

	p := fromString(s)
	s2 := toString(p, false)

	if s != s2 {
		t.Error("Strings do not match")
	}
}
