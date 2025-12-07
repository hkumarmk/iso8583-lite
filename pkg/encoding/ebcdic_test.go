package encoding

import (
	"bytes"
	"testing"
)

func TestEBCDIC037_EncodeDecode(t *testing.T) {
	ascii := []byte("0123456789ABCDEFabcdef")

	enc, err := EBCDIC037.Encode(ascii)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	dec, n, err := EBCDIC037.Decode(enc)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if n != len(enc) {
		t.Errorf("Decode did not consume all input: got %d, want %d", n, len(enc))
	}

	if !bytes.Equal(ascii, dec) {
		t.Errorf("Round-trip Encode/Decode failed.\nInput:  %v\nOutput: %v", ascii, dec)
	}
}

func TestEBCDIC037_Encode_InvalidASCII(t *testing.T) {
	in := []byte{0x80, 0xFF}

	_, err := EBCDIC037.Encode(in)
	if err == nil {
		t.Error("expected error for non-ASCII input, got nil")
	}
}

// Only test round-trip for safe EBCDIC ASCII subset: A-Z, 0-9, space, and a few common symbols.
func TestEBCDIC037_EncodeDecode_SafeSubset(t *testing.T) {
	safe := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 -./")

	enc, err := EBCDIC037.Encode(safe)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	dec, n, err := EBCDIC037.Decode(enc)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if n != len(enc) {
		t.Errorf("Decode did not consume all input: got %d, want %d", n, len(enc))
	}

	if !bytes.Equal(safe, dec) {
		t.Errorf("Safe subset round-trip failed.\nInput:  %v\nOutput: %v", safe, dec)
	}
}

func TestEBCDIC037_EncodeDecode_ControlChars(t *testing.T) {
	for b := byte(0x00); b <= 0x1F; b++ {
		in := []byte{b}

		enc, err := EBCDIC037.Encode(in)
		if err != nil {
			t.Errorf("Encode failed for control char 0x%02X: %v", b, err)

			continue
		}

		dec, n, err := EBCDIC037.Decode(enc)
		if err != nil {
			t.Errorf("Decode failed for control char 0x%02X: %v", b, err)

			continue
		}

		if n != len(enc) {
			t.Errorf("Decode did not consume all input for 0x%02X: got %d, want %d", b, n, len(enc))
		}

		if !bytes.Equal(in, dec) {
			t.Errorf("Control char round-trip failed for 0x%02X: got %v", b, dec)
		}
	}
	// DEL (0x7F)
	in := []byte{0x7F}

	enc, err := EBCDIC037.Encode(in)
	if err != nil {
		t.Errorf("Encode failed for DEL: %v", err)
	}

	dec, n, err := EBCDIC037.Decode(enc)
	if err != nil {
		t.Errorf("Decode failed for DEL: %v", err)
	}

	if n != len(enc) {
		t.Errorf("Decode did not consume all input for DEL: got %d, want %d", n, len(enc))
	}

	if !bytes.Equal(in, dec) {
		t.Errorf("DEL round-trip failed: got %v", dec)
	}
}

func TestEBCDIC037_EncodeDecode_EmptyAndSingleByte(t *testing.T) {
	cases := [][]byte{
		{},
		{0x00},
		{0x41},
		{0x7F},
	}
	for _, in := range cases {
		enc, err := EBCDIC037.Encode(in)
		if err != nil {
			t.Errorf("Encode failed for %v: %v", in, err)

			continue
		}

		dec, n, err := EBCDIC037.Decode(enc)
		if err != nil {
			t.Errorf("Decode failed for %v: %v", in, err)

			continue
		}

		if n != len(enc) {
			t.Errorf("Decode did not consume all input for %v: got %d, want %d", in, n, len(enc))
		}

		if !bytes.Equal(in, dec) {
			t.Errorf("Round-trip failed for %v: got %v", in, dec)
		}
	}
}

func TestEBCDIC037_KnownPairs(t *testing.T) {
	pairs := []struct {
		ascii  string
		ebcdic []byte
	}{
		{"HELLO", []byte{0xC8, 0xC5, 0xD3, 0xD3, 0xD6}},
		{"1234", []byte{0xF1, 0xF2, 0xF3, 0xF4}},
		{"!@#", []byte{0x5A, 0x7C, 0x7B}},
	}
	for _, p := range pairs {
		enc, err := EBCDIC037.Encode([]byte(p.ascii))
		if err != nil {
			t.Errorf("Encode failed for %q: %v", p.ascii, err)

			continue
		}

		if !bytes.Equal(enc, p.ebcdic) {
			t.Errorf("Encode mismatch for %q: got %v, want %v", p.ascii, enc, p.ebcdic)
		}

		dec, n, err := EBCDIC037.Decode(p.ebcdic)
		if err != nil {
			t.Errorf("Decode failed for %v: %v", p.ebcdic, err)

			continue
		}

		if n != len(p.ebcdic) {
			t.Errorf("Decode did not consume all input for %v: got %d, want %d", p.ebcdic, n, len(p.ebcdic))
		}

		if !bytes.Equal(dec, []byte(p.ascii)) {
			t.Errorf("Decode mismatch for %v: got %v, want %v", p.ebcdic, dec, p.ascii)
		}
	}
}

func TestEBCDIC037_ISO8583FieldValues(t *testing.T) {
	tests := []struct {
		name  string
		ascii []byte
	}{
		{"Digits", []byte("0123456789")},
		{"Uppercase Letters", []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")},
		{"PAN Separators", []byte("0123=4567")},
		{"Currency Code", []byte("USD")},
		{"Account Number", []byte("1234567890123456")},
		{"Symbols", []byte(" -./")},
		{"Empty", []byte("")},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			enc, err := EBCDIC037.Encode(tc.ascii)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			dec, n, err := EBCDIC037.Decode(enc)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}

			if n != len(enc) {
				t.Errorf("Decode did not consume all input for %q: got %d, want %d", tc.ascii, n, len(enc))
			}

			if !bytes.Equal(tc.ascii, dec) {
				t.Errorf("Round-trip failed for %q: got %v", tc.ascii, dec)
			}
		})
	}
}
