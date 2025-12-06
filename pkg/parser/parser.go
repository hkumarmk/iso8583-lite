package parser

import (
	"fmt"

	"github.com/hkumarmk/iso8583-lite/pkg/spec"
)

// Parser is a stateless field location calculator that uses a Spec.
// It calculates where fields are located in a message buffer without storing state.
type Parser struct {
	spec *spec.Spec
}

// Bitmap interface for bitmap operations (to avoid circular dependency).
type Bitmap interface {
	IsSet(fieldNum int) bool
	PresentFields() []int
	IsExtended() bool
}

// NewParser creates a new stateless parser for the given spec.
func NewParser(s *spec.Spec) *Parser {
	return &Parser{
		spec: s,
	}
}

// ParseField calculates the cursor for a field based on the spec.
// Requires the buffer, field number, and starting offset.
// Returns cursor and error if field cannot be parsed.
func (p *Parser) ParseField(buf []byte, fieldNum int, offset int) (Cursor, error) {
	fieldSpec, ok := p.spec.Fields[fieldNum]
	if !ok {
		return Cursor{}, fmt.Errorf("field %d not defined in spec", fieldNum)
	}

	// Check if we have enough data
	if offset >= len(buf) {
		return Cursor{}, fmt.Errorf("field %d: offset %d exceeds buffer length %d", fieldNum, offset, len(buf))
	}

	// Parse based on field type
	switch fieldSpec.Type {
	case spec.FieldTypeFixed:
		return p.parseFixed(buf, fieldSpec, offset)
	case spec.FieldTypeL, spec.FieldTypeLL, spec.FieldTypeLLL:
		return p.parseVariable(buf, fieldSpec, offset)
	case spec.FieldTypeBitmap:
		return p.parseBitmap(buf, fieldSpec, offset)
	default:
		return Cursor{}, fmt.Errorf("unsupported field type %v", fieldSpec.Type)
	}
}

// parseFixed parses a fixed-length field.
func (p *Parser) parseFixed(buf []byte, fieldSpec *spec.FieldSpec, offset int) (Cursor, error) {
	if offset+fieldSpec.Length > len(buf) {
		return Cursor{}, fmt.Errorf("field %d (%s): expected %d bytes at offset %d, buffer has %d bytes",
			fieldSpec.Number, fieldSpec.Name, fieldSpec.Length, offset, len(buf))
	}

	return Cursor{
		Start: offset,
		End:   offset + fieldSpec.Length,
	}, nil
}

// parseVariable parses a variable-length field (L, LL, LLL).
func (p *Parser) parseVariable(buf []byte, fieldSpec *spec.FieldSpec, offset int) (Cursor, error) {
	lenDigits := fieldSpec.Type.LengthIndicatorDigits()

	// Check if we have enough bytes for the length indicator
	if offset+lenDigits > len(buf) {
		return Cursor{}, fmt.Errorf("field %d (%s): expected %d bytes for length indicator at offset %d, buffer has %d bytes",
			fieldSpec.Number, fieldSpec.Name, lenDigits, offset, len(buf))
	}

	// Parse length indicator
	lenBytes := buf[offset : offset+lenDigits]
	fieldLen, err := parseInt(lenBytes)
	if err != nil {
		return Cursor{}, fmt.Errorf("field %d (%s): invalid length indicator %q: %w",
			fieldSpec.Number, fieldSpec.Name, string(lenBytes), err)
	}

	// Validate field length
	if fieldLen > fieldSpec.MaxLength {
		return Cursor{}, fmt.Errorf("field %d (%s): length %d exceeds max length %d",
			fieldSpec.Number, fieldSpec.Name, fieldLen, fieldSpec.MaxLength)
	}

	dataStart := offset + lenDigits
	dataEnd := dataStart + fieldLen

	// Check if we have enough data
	if dataEnd > len(buf) {
		return Cursor{}, fmt.Errorf("field %d (%s): expected %d bytes of data at offset %d, buffer has %d bytes",
			fieldSpec.Number, fieldSpec.Name, fieldLen, dataStart, len(buf))
	}

	return Cursor{
		Start: dataStart,
		End:   dataEnd,
	}, nil
}

// parseBitmap parses a bitmap field (8 or 16 bytes).
func (p *Parser) parseBitmap(buf []byte, fieldSpec *spec.FieldSpec, offset int) (Cursor, error) {
	if offset+fieldSpec.Length > len(buf) {
		return Cursor{}, fmt.Errorf("field %d (%s): expected %d bytes at offset %d, buffer has %d bytes",
			fieldSpec.Number, fieldSpec.Name, fieldSpec.Length, offset, len(buf))
	}

	return Cursor{
		Start: offset,
		End:   offset + fieldSpec.Length,
	}, nil
}

// parseInt parses a numeric byte slice into an integer.
func parseInt(b []byte) (int, error) {
	result := 0
	for _, c := range b {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid digit %q", c)
		}
		result = result*10 + int(c-'0')
	}
	return result, nil
}
