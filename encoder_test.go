package xdelta

import (
	"testing"
)

func testInitEncoderBase(t *testing.T, options EncoderOptions) {
	enc, err := NewEncoder(options)
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

func TestEncoderInitEmpty(t *testing.T) {
	testInitEncoderBase(t, EncoderOptions{})
}

func TestEncoderInitFull(t *testing.T) {
	testInitEncoderBase(t, EncoderOptions{
		FileID:      "test.file",
		BlockSizeKB: 4096,
		Header:      []byte(`{"myheader":true}`),
	})
}
