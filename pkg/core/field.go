package core

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/hkumarmk/iso8583-lite/pkg/parser"
	"github.com/hkumarmk/iso8583-lite/pkg/spec"
)

// FieldAccessor provides zero-copy field access with type conversion.
type FieldAccessor interface {
	// Exists returns true if the field is present in the message.
	Exists() bool

	// Bytes returns the raw field bytes (zero-copy slice).
	Bytes() []byte

	// String returns the field value as a string.
	String() string

	// Int returns the field value as an int (returns 0 on error).
	Int() int

	// IntE returns the field value as an int with error handling.
	IntE() (int, error)

	// Int64 returns the field value as an int64 (returns 0 on error).
	Int64() int64

	// Int64E returns the field value as an int64 with error handling.
	Int64E() (int64, error)

	// Hex returns the field value as a hex string.
	Hex() string

	// Len returns the length of the field in bytes.
	Len() int
}

// Field provides access to a single ISO8583 field with multiple type accessors.
// Supports hierarchical structure for composite fields with lazy parsing and caching.
type Field struct {
	data     []byte
	exists   bool
	spec     *spec.FieldSpec // Spec for this field (defines structure, children)
	parser   *parser.Parser  // Parser for lazy parsing children
	children map[int]*Field  // Subfields (lazy-loaded, nil until first access)
}

var _ FieldAccessor = (*Field)(nil)

// NewField creates a field with raw data and existence flag.
func NewField(data []byte, exists bool) *Field {
	return &Field{
		data:   data,
		exists: exists,
	}
}

// NewFieldWithSpec creates a field with a spec and parser for lazy child parsing.
func NewFieldWithSpec(data []byte, exists bool, fieldSpec *spec.FieldSpec, p *parser.Parser) *Field {
	return &Field{
		data:   data,
		exists: exists,
		spec:   fieldSpec,
		parser: p,
	}
}

// Exists returns true if the field is present.
func (f *Field) Exists() bool {
	return f.exists
}

// Bytes returns the raw field data if present, or nil.
func (f *Field) Bytes() []byte {
	if !f.exists {
		return nil
	}

	return f.data
}

// String returns the field data as a string, or an empty string if not present.
func (f *Field) String() string {
	if !f.exists {
		return ""
	}

	return string(f.data)
}

// Int returns the field value as int, or zero if not present or invalid.
func (f *Field) Int() int {
	val, _ := f.IntE()

	return val
}

// IntE returns the field value as int, or an error if not present or invalid.
func (f *Field) IntE() (int, error) {
	if !f.exists {
		return 0, ErrFieldNotPresent
	}

	val, err := strconv.Atoi(f.String())
	if err != nil {
		return 0, fmt.Errorf("failed to convert field to int: %w", err)
	}

	return val, nil
}

// Int64 returns the field value as int64, or zero if not present or invalid.
func (f *Field) Int64() int64 {
	val, _ := f.Int64E()

	return val
}

// Int64E returns the field value as int64, or an error if not present or invalid.
func (f *Field) Int64E() (int64, error) {
	if !f.exists {
		return 0, ErrFieldNotPresent
	}

	val, err := strconv.ParseInt(f.String(), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert field to int64: %w", err)
	}

	return val, nil
}

// Hex returns the field data as a hex string, or an empty string if not present.
func (f *Field) Hex() string {
	if !f.exists {
		return ""
	}

	return hex.EncodeToString(f.data)
}

// Len returns the length of the field data in bytes.
func (f *Field) Len() int {
	return len(f.data)
}

// Subfield returns a child field by number for composite fields. Returns a non-existent field if not found.
func (f *Field) Subfield(num int) *Field {
	if !f.exists {
		return &Field{exists: false}
	}

	// Check cache first
	if f.children != nil {
		if child, ok := f.children[num]; ok {
			return child
		}
	}

	// Not cached - try to parse if we have spec and parser
	if f.spec != nil && f.parser != nil && len(f.spec.Children) > 0 {
		// Find child spec
		var childSpec *spec.FieldSpec

		for _, cs := range f.spec.Children {
			if cs.Number == num {
				childSpec = cs

				break
			}
		}

		if childSpec != nil {
			// Calculate offset within this field's data
			// TODO: Need to walk through previous siblings to calculate offset
			// For now, return non-existent field
			// This will be implemented when we integrate with Message
			return &Field{exists: false}
		}
	}

	// Not found or can't parse
	return &Field{exists: false}
}

// SetSubfield sets a child field for this field.
func (f *Field) SetSubfield(num int, child *Field) {
	if f.children == nil {
		f.children = make(map[int]*Field)
	}

	f.children[num] = child
}

// HasSubfields returns true if this field has parsed subfields.
func (f *Field) HasSubfields() bool {
	return len(f.children) > 0
}
