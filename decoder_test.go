package xdelta

import (
	"testing"
)

func testInitDecoderBase(t *testing.T, options DecoderOptions) {
	enc, err := NewDecoder(options)
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	if enc.handle == nil {
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
	if enc.handle != nil {
		t.Errorf("Handle is expected to be 0 after close")
	}

	// multiple calls may not fail
	err = enc.Close()
	if err != nil {
		t.Errorf("Failed to close: %v", err)
	}
}

func TestDecoderInitEmpty(t *testing.T) {
	testInitDecoderBase(t, DecoderOptions{})
}

func TestDecoderInitFull(t *testing.T) {
	testInitDecoderBase(t, DecoderOptions{
		FileID:      "test.file",
		BlockSizeKB: 4096,
	})
}
