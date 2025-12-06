package parser

import (
	"testing"

	"github.com/hkumarmk/iso8583-lite/pkg/spec"
)

func TestParseFixed(t *testing.T) {
	s := &spec.Spec{
		Fields: map[int]*spec.FieldSpec{
			4: {
				Number:   4,
				Type:     spec.FieldTypeFixed,
				Length:   12,
				DataType: spec.DataTypeNumeric,
			},
		},
	}

	buf := []byte("0100000000000001000000001000")
	p := NewParser(s)

	cur, err := p.ParseField(buf, 4, 16)
	if err != nil {
		t.Fatalf("ParseField() error = %v", err)
	}

	if cur.Start != 16 || cur.End != 28 {
		t.Errorf("ParseField() cursor = {%d, %d}, want {16, 28}", cur.Start, cur.End)
	}

	data := cur.Extract(buf)
	if string(data) != "000000001000" {
		t.Errorf("extracted data = %q, want %q", string(data), "000000001000")
	}
}

func TestParseFixedTooShort(t *testing.T) {
	s := &spec.Spec{
		Fields: map[int]*spec.FieldSpec{
			4: {
				Number:   4,
				Type:     spec.FieldTypeFixed,
				Length:   12,
				DataType: spec.DataTypeNumeric,
			},
		},
	}

	buf := []byte("0100")
	p := NewParser(s)

	_, err := p.ParseField(buf, 4, 0)
	if err == nil {
		t.Error("ParseField() expected error for short buffer, got nil")
	}
}

func TestParseVariable(t *testing.T) {
	s := &spec.Spec{
		Fields: map[int]*spec.FieldSpec{
			2: {
				Number:    2,
				Type:      spec.FieldTypeLL,
				MaxLength: 19,
				DataType:  spec.DataTypeNumeric,
			},
		},
	}

	// "16" = length, "1234567890123456" = PAN
	buf := []byte("161234567890123456")
	p := NewParser(s)

	cur, err := p.ParseField(buf, 2, 0)
	if err != nil {
		t.Fatalf("ParseField() error = %v", err)
	}

	// Cursor should point to data (after length indicator)
	if cur.Start != 2 || cur.End != 18 {
		t.Errorf("ParseField() cursor = {%d, %d}, want {2, 18}", cur.Start, cur.End)
	}

	data := cur.Extract(buf)
	if string(data) != "1234567890123456" {
		t.Errorf("extracted data = %q, want %q", string(data), "1234567890123456")
	}
}

func TestParseVariableInvalidLength(t *testing.T) {
	s := &spec.Spec{
		Fields: map[int]*spec.FieldSpec{
			2: {
				Number:    2,
				Type:      spec.FieldTypeLL,
				MaxLength: 19,
				DataType:  spec.DataTypeNumeric,
			},
		},
	}

	buf := []byte("XX1234567890123456")
	p := NewParser(s)

	_, err := p.ParseField(buf, 2, 0)
	if err == nil {
		t.Error("ParseField() expected error for invalid length indicator, got nil")
	}
}

func TestParseVariableExceedsMaxLength(t *testing.T) {
	s := &spec.Spec{
		Fields: map[int]*spec.FieldSpec{
			2: {
				Number:    2,
				Type:      spec.FieldTypeLL,
				MaxLength: 10,
				DataType:  spec.DataTypeNumeric,
			},
		},
	}

	buf := []byte("161234567890123456")
	p := NewParser(s)

	_, err := p.ParseField(buf, 2, 0)
	if err == nil {
		t.Error("ParseField() expected error for length exceeding max, got nil")
	}
}

func TestParseVariableTruncatedData(t *testing.T) {
	s := &spec.Spec{
		Fields: map[int]*spec.FieldSpec{
			2: {
				Number:    2,
				Type:      spec.FieldTypeLL,
				MaxLength: 19,
				DataType:  spec.DataTypeNumeric,
			},
		},
	}

	buf := []byte("16123456") // Says 16 bytes but only has 6
	p := NewParser(s)

	_, err := p.ParseField(buf, 2, 0)
	if err == nil {
		t.Error("ParseField() expected error for truncated data, got nil")
	}
}

func TestParseFieldUnsupportedType(t *testing.T) {
	s := &spec.Spec{
		Fields: map[int]*spec.FieldSpec{
			99: {
				Number: 99,
				Type:   spec.FieldType(999), // Invalid type
			},
		},
	}

	buf := []byte("test")
	p := NewParser(s)

	_, err := p.ParseField(buf, 99, 0)
	if err == nil {
		t.Error("ParseField() expected error for unsupported type, got nil")
	}
}

func TestParseFieldNotInSpec(t *testing.T) {
	s := &spec.Spec{
		Fields: map[int]*spec.FieldSpec{},
	}

	buf := []byte("test")
	p := NewParser(s)

	_, err := p.ParseField(buf, 42, 0)
	if err == nil {
		t.Error("ParseField() expected error for field not in spec, got nil")
	}
}
