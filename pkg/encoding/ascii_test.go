package encoding

import (
	"bytes"
	"testing"
)

func TestASCII_EncodeDecode_Valid(t *testing.T) {
	input := []byte("0123456789ABCDefghijklmnopqrstuvwxyz!@#$%^&*()_+-=")

	enc, err := ASCII.Encode(input)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if !bytes.Equal(enc, input) {
		t.Errorf("Encode should be a no-op for ASCII: got %v, want %v", enc, input)
	}

	dec, read, err := ASCII.Decode(enc)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if read != len(enc) {
		t.Errorf("Decode: n = %d, want %d", read, len(enc))
	}

	if !bytes.Equal(dec, input) {
		t.Errorf("Decode should be a no-op for ASCII: got %v, want %v", dec, input)
	}
}

func TestASCII_Encode_NonASCII(t *testing.T) {
	input := []byte{0x41, 0x80, 0xFF}

	_, err := ASCII.Encode(input)
	if err == nil {
		t.Error("expected error for non-ASCII input, got nil")
	}
}

func TestASCII_Decode_NonASCII(t *testing.T) {
	input := []byte{0x41, 0x80, 0xFF}

	_, read, err := ASCII.Decode(input)
	if err == nil {
		t.Error("expected error for non-ASCII input, got nil")
	}

	if read != 0 {
		t.Errorf("Decode: n = %d, want 0 for error", read)
	}
}

func TestASCII_Name(t *testing.T) {
	if ASCII.Name() != "ASCII" {
		t.Errorf("Name() = %q, want %q", ASCII.Name(), "ASCII")
	}
}
