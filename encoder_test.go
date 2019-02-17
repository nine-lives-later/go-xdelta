package xdelta

import (
	"testing"

	"github.com/konsorten/go-xdelta/xdelta-lib"
)

func testInitBase(t *testing.T, options EncoderOptions) {
	enc, err := NewEncoder(options)
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	if enc.handle == 0 {
		t.Fatalf("Handle is expected to be valid after init")
	}

	err = enc.GetStreamError()
	if err != nil {
		t.Errorf("Stream returned error: %v", err)
	}

	err = enc.Close()
	if err != nil {
		t.Errorf("Failed to close: %v", err)
	}
	if enc.handle != 0 {
		t.Errorf("Handle is expected to be 0 after close")
	}

	// multiple calls may not fail
	err = enc.Close()
	if err != nil {
		t.Errorf("Failed to close: %v", err)
	}
}

func TestInitEmpty(t *testing.T) {
	testInitBase(t, EncoderOptions{})
}

func TestInitFull(t *testing.T) {
	testInitBase(t, EncoderOptions{
		FileID:      "test.file",
		BlockSizeKB: 4096,
		Header:      []byte(`{"myheader":true}`),
	})
}

func TestInitInvalidBlockSize(t *testing.T) {
	_, err := NewEncoder(EncoderOptions{BlockSizeKB: -1})
	if err != nil {
		if errX, ok := err.(lib.XdeltaError); ok {
			if errX == lib.XdeltaError_ArgumentOutOfRange {
				return
			}
		}

		t.Fatalf("Wrong error: %v", err)
	}

	t.Fatalf("Expected error")
}
