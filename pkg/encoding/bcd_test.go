package encoding

import (
	"bytes"
	"testing"
)

func TestBCD_EncodeDecode(t *testing.T) {
	cases := []struct {
		name  string
		ascii string
		bcd   []byte
	}{
		{"Even digits", "1234", []byte{0x12, 0x34}},
		{"Odd digits", "123", []byte{0x01, 0x23}},
		{"Single digit", "7", []byte{0x07}},
		{"Empty", "", []byte{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			enc, err := BCD.Encode([]byte(tc.ascii))
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			if !bytes.Equal(enc, tc.bcd) {
				t.Errorf("Encode: got %v, want %v", enc, tc.bcd)
			}

			dec, n, err := BCD.Decode(enc)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}

			if n != len(enc) {
				t.Errorf("Decode did not consume all input: got %d, want %d", n, len(enc))
			}

			want := tc.ascii
			if len(want)%2 != 0 {
				want = "0" + want
			}

			if string(dec) != want {
				t.Errorf("Decode: got %q, want %q", dec, want)
			}
		})
	}
}

func TestBCD_Encode_Invalid(t *testing.T) {
	_, err := BCD.Encode([]byte("12A4"))
	if err == nil {
		t.Error("expected error for non-digit input, got nil")
	}
}

func TestBCD_Decode_Invalid(t *testing.T) {
	_, _, err := BCD.Decode([]byte{0x1A})
	if err == nil {
		t.Error("expected error for invalid BCD digit, got nil")
	}
}

func TestBCD_Name(t *testing.T) {
	if BCD.Name() != "BCD" {
		t.Errorf("Name() = %q, want %q", BCD.Name(), "BCD")
	}
}
