package lib

import (
	"testing"
)

func TestReturnErrorNotImplemented(t *testing.T) {
	err := callToError(testReturnErrorNotImplemented.Call())
	if err != nil {
		if err == XdeltaError_NotImplemented {
			return
		}
		t.Fatalf("Got wrong error: %v", err)
	}
	t.Fatal("Expected error")
}
