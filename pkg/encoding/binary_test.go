package encoding

import (
	"bytes"
	"testing"
)

func TestBinaryEncoder(t *testing.T) {
	cases := []struct {
		name string
		in   []byte
	}{
		{"empty", []byte{}},
		{"ascii", []byte("hello")},
		{"binary", []byte{0x00, 0xFF, 0x10, 0x20}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			enc, err := Binary.Encode(tc.in)
			if err != nil {
				t.Fatalf("Encode error: %v", err)
			}
			if !bytes.Equal(enc, tc.in) {
				t.Errorf("Encode mismatch: got %v, want %v", enc, tc.in)
			}
			dec, n, err := Binary.Decode(enc)
			if err != nil {
				t.Fatalf("Decode error: %v", err)
			}
			if n != len(enc) {
				t.Errorf("Decode did not consume all input: got %d, want %d", n, len(enc))
			}
			if !bytes.Equal(dec, tc.in) {
				t.Errorf("Decode mismatch: got %v, want %v", dec, tc.in)
			}
		})
	}
}
