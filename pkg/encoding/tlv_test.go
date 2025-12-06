package encoding

import (
	"bytes"
	"testing"

	"github.com/euicc-go/bertlv"
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
	tag := bertlv.NewTag(bertlv.Application, bertlv.Primitive, 0x5F2A)
	value := []byte{0x01, 0x02}
	tlv := bertlv.NewValue(tag, value)
	data, err := tlv.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}

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

func TestTLVEncoder_EncodeDecode_MultipleTLVs(t *testing.T) {
	tag1 := bertlv.NewTag(bertlv.Application, bertlv.Primitive, 0x9F33)
	val1 := []byte{0x01, 0x02, 0x03}
	tlv1 := bertlv.NewValue(tag1, val1)

	tag2 := bertlv.NewTag(bertlv.Application, bertlv.Primitive, 0x95)
	val2 := []byte{0xAA, 0xBB}
	tlv2 := bertlv.NewValue(tag2, val2)

	data1, err := tlv1.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	data2, err := tlv2.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	data := append(data1, data2...)

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

func TestTLVEncoder_NestedTLV(t *testing.T) {
	// Construct a TLV with a constructed tag and a child TLV
	tagOuter := bertlv.NewTag(bertlv.Application, bertlv.Constructed, 0xE1)
	tagInner := bertlv.NewTag(bertlv.Application, bertlv.Primitive, 0x5F2A)
	valInner := []byte{0x01, 0x02}
	inner := bertlv.NewValue(tagInner, valInner)
	outer := bertlv.NewChildren(tagOuter, inner)
	data, err := outer.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
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

func TestTLVEncoder_LongFormLength(t *testing.T) {
	// TLV with long-form length (length > 127)
	tag := bertlv.NewTag(bertlv.Application, bertlv.Primitive, 0x5F2A)
	val := bytes.Repeat([]byte{0xAB}, 130)
	tlv := bertlv.NewValue(tag, val)
	data, err := tlv.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
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

func TestTLVEncoder_ZeroLengthValue(t *testing.T) {
	tag := bertlv.NewTag(bertlv.Application, bertlv.Primitive, 0x5F2A)
	tlv := bertlv.NewValue(tag, []byte{})
	data, err := tlv.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
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
