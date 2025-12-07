package encoding

import (
	"bytes"
	"testing"
)

func TestTLVEncoder_InvalidInput(t *testing.T) {
	cases := [][]byte{
		{},                 // empty
		{0x00},             // incomplete tag
		{0x9F},             // incomplete tag
		{0x9F, 0x33},       // tag but no length/value
		{0x9F, 0x33, 0xFF}, // tag + invalid length
		{0x9F, 0x33, 0x01}, // tag + length, but missing value
	}
	for i, in := range cases {
		_, _, err := TLV.Decode(in)
		if len(in) == 0 {
			// Empty input is valid, should not error
			if err != nil {
				t.Errorf("case %d: expected no error for empty input, got %v", i, err)
			}
		} else {
			if err == nil {
				t.Errorf("case %d: expected error for invalid input %v, got nil", i, in)
			}
		}
	}
}

func TestTLVEncoder_EncodeDecode_Simple(t *testing.T) {
	// Flat TLV: tag=0x5F, length=2, value=0x01 0x02
	data := []byte{0x5F, 0x02, 0x01, 0x02}

	enc, err := TLV.Encode(data)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if !bytes.Equal(enc, data) {
		t.Errorf("Encode mismatch: got %v, want %v", enc, data)
	}

	dec, n, err := TLV.Decode(enc)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if n != len(enc) {
		t.Errorf("Decode did not consume all input: got %d, want %d", n, len(enc))
	}

	if !bytes.Equal(dec, data) {
		t.Errorf("Decode mismatch: got %v, want %v", dec, data)
	}
}

func TestTLVEncoder_TableDriven(t *testing.T) {
	cases := []struct {
		name string
		data []byte
	}{
		{
			name: "DE55 sample-1",
			data: []byte{
				0x95, 0x05, 0x80, 0x00, 0x00, 0x00, 0x00,
				0x9A, 0x03, 0x20, 0x12, 0x31,
				0x5F, 0x2A, 0x02, 0x08, 0x40,
				0x9F, 0x02, 0x06, 0x00, 0x00, 0x01, 0x00, 0x00,
				0x9F, 0x03, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x9F, 0x1A, 0x02, 0x08, 0x40,
				0x9F, 0x27, 0x01, 0x80,
				0x9F, 0x36, 0x02, 0x00, 0x3C,
				0x9F, 0x37, 0x04, 0x6B, 0x1A, 0x2C, 0x3D,
			},
		},
		{
			name: "DE55 sample-2",
			data: []byte{
				0x95, 0x05, 0x00, 0x00, 0x80, 0x00, 0x00,
				0x9A, 0x03, 0x24, 0x10, 0x03,
				0x5F, 0x2A, 0x02, 0x07, 0x10,
				0x9F, 0x02, 0x06, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00,
				0x9F, 0x10, 0x12, 0x01, 0x10, 0x20, 0x80, 0x03, 0x24, 0x20, 0x00,
				0x96, 0x1F, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF,
				0x9F, 0x1A, 0x02, 0x07, 0x10,
				0x9F, 0x26, 0x08, 0xA3, 0xFD, 0xE2, 0xBF, 0x27, 0xF3, 0x98, 0x39,
				0x9F, 0x27, 0x01, 0x00,
				0x9F, 0x36, 0x02, 0x00, 0x19,
				0x9F, 0x37, 0x04, 0x1C, 0x2D, 0x1D, 0xBE,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			enc, err := TLV.Encode(tc.data)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			if !bytes.Equal(enc, tc.data) {
				t.Errorf("Encode mismatch: got %v, want %v", enc, tc.data)
			}

			dec, n, err := TLV.Decode(enc)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}

			if n != len(enc) {
				t.Errorf("Decode did not consume all input: got %d, want %d", n, len(enc))
			}

			if !bytes.Equal(dec, tc.data) {
				t.Errorf("Decode mismatch: got %v, want %v", dec, tc.data)
			}
		})
	}
}

func TestTLVEncoder_Empty(t *testing.T) {
	enc, err := TLV.Encode([]byte{})
	if err != nil {
		t.Fatalf("Encode failed for empty: %v", err)
	}

	if len(enc) != 0 {
		t.Errorf("Encode: expected empty, got %v", enc)
	}

	dec, n, err := TLV.Decode(enc)
	if err != nil {
		t.Fatalf("Decode failed for empty: %v", err)
	}

	if n != 0 {
		t.Errorf("Decode: expected 0 bytes read, got %d", n)
	}

	if len(dec) != 0 {
		t.Errorf("Decode: expected empty, got %v", dec)
	}
}
