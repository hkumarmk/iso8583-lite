package encoding

import (
	"bytes"
	"testing"
)

func TestHexEncoder(t *testing.T) {
	cases := []struct {
		name string
		in   []byte
		out  string
	}{
		{"empty", []byte{}, ""},
		{"ascii", []byte("hi"), "6869"},
		{"binary", []byte{0x00, 0xFF, 0x10}, "00ff10"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			enc, err := Hex.Encode(tc.in)
			if err != nil {
				t.Fatalf("Encode error: %v", err)
			}
			if string(enc) != tc.out {
				t.Errorf("Encode mismatch: got %q, want %q", enc, tc.out)
			}
			dec, n, err := Hex.Decode(enc)
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
