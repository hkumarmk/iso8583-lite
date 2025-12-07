package core

import (
	"fmt"

	"github.com/hkumarmk/iso8583-lite/pkg/parser"
	"github.com/hkumarmk/iso8583-lite/pkg/spec"
)

// MessageReader defines the complete interface for reading and validating ISO8583 messages.
type MessageReader interface {
	// Parse parses the MTI, bitmap, and all present fields (eager parsing).
	// For full structural validation, use Validate() with StructuralValidator.
	Parse() error

	// MTI returns the Message Type Indicator field.
	MTI() FieldAccessor

	// Field returns the accessor for the specified field number.
	Field(fieldNum int) FieldAccessor

	// HasField returns true if the field is present.
	HasField(fieldNum int) bool

	// PresentFields returns all present field numbers.
	PresentFields() []int

	// Bytes returns the raw message bytes.
	Bytes() []byte

	// Validate performs validation using the provided validator.
	// Pass nil to skip validation. Use NewCompositeValidator() to combine validators.
	Validate(validator Validator) error

	// ValidateField validates a specific field.
	ValidateField(fieldNum int) error
}

// mtiLength is the length of the Message Type Indicator field in bytes.
const (
	mtiLength        = 4
	minMessageLength = 12 // Minimum length: MTI (4) + Primary Bitmap (8)
)

// Message represents a parsed ISO8583 message with zero-copy field access.
type Message struct {
	buf     []byte                // Raw message bytes
	bitmap  *Bitmap               // Bitmap indicating present fields
	mti     string                // Message Type Indicator (field 0)
	cursors map[int]parser.Cursor // Cached field positions (eager parsing)
	spec    *spec.Spec            // Message specification
	parser  *parser.Parser        // Parser for field parsing
}

var _ MessageReader = (*Message)(nil)

// NewMessage creates a message wrapper with the given spec.
// Use Parse() to parse MTI, bitmap, and all fields.
func NewMessage(buf []byte, s *spec.Spec) *Message {
	return &Message{
		buf:     buf,
		cursors: make(map[int]parser.Cursor),
		spec:    s,
		parser:  parser.NewParser(s),
	}
}

// Parse parses the MTI, bitmap, and all present fields in the ISO8583 message.
// It validates the MTI structure and bitmap, and performs eager parsing of all fields.
func (m *Message) Parse() error {
	if len(m.buf) < mtiLength {
		return ErrInvalidMTI(len(m.buf))
	}

	m.mti = string(m.buf[0:mtiLength])

	if !isValidMTIStructure(m.mti) {
		return ErrInvalidMTIFormat(m.mti)
	}

	if len(m.buf) < minMessageLength {
		return ErrMessageTooShort(minMessageLength, len(m.buf))
	}

	bitmap, _, err := NewBitmap(m.buf[4:])
	if err != nil {
		return ErrBitmapParseFailed(err)
	}

	m.bitmap = bitmap

	// Eager parsing: Parse all present fields immediately
	if err := m.parseAllFields(); err != nil {
		return err
	}

	return nil
}

// MTI returns the Message Type Indicator field accessor.
//
//nolint:ireturn // Returning interface for extensibility is intentional
func (m *Message) MTI() FieldAccessor {
	return NewField([]byte(m.mti), m.mti != "")
}

// Field returns the accessor for the specified field number.
//
//nolint:ireturn // Returning interface for extensibility is intentional
func (m *Message) Field(fieldNum int) FieldAccessor {
	if fieldNum < 0 || fieldNum > 128 {
		return NewField(nil, false)
	}

	if fieldNum == 0 {
		return m.MTI()
	}

	if m.bitmap == nil || !m.bitmap.IsSet(fieldNum) {
		return NewField(nil, false)
	}

	// Get cursor from cache (populated during Parse())
	cursor, ok := m.cursors[fieldNum]
	if !ok {
		// Field present in bitmap but not parsed (shouldn't happen after Parse())
		return NewField(nil, false)
	}

	// Extract data using cursor (zero-copy)
	data := cursor.Extract(m.buf)
	if data == nil {
		return NewField(nil, false)
	}

	// Get field spec for this field
	fieldSpec := m.spec.Fields[fieldNum]

	// Create field with spec and parser for subfield access
	return NewFieldWithSpec(data, true, fieldSpec, m.parser)
}

// HasField returns true if the specified field is present in the message.
func (m *Message) HasField(fieldNum int) bool {
	if fieldNum == 0 {
		return true
	}

	if m.bitmap == nil {
		return false
	}

	return m.bitmap.IsSet(fieldNum)
}

// Bytes returns the raw ISO8583 message bytes.
func (m *Message) Bytes() []byte {
	return m.buf
}

// PresentFields returns all present field numbers, including the MTI (field 0).
func (m *Message) PresentFields() []int {
	if m.bitmap == nil {
		return []int{0}
	}

	fields := []int{0} // Start with MTI
	fields = append(fields, m.bitmap.PresentFields()...)

	return fields
}

// Validate performs validation using the provided validator.
// Pass nil to skip validation.
// Use NewCompositeValidator() to combine multiple validators.
// Example:
//
//	validator := NewCompositeValidator(
//	    NewStructuralValidator(spec),
//	    NewFormatValidator(spec),
//	    NewBusinessValidator(spec, rules...),
//	)
//	if err := msg.Validate(validator); err != nil {
//	    // Handle validation error
//	}
//
//nolint:wrapcheck // Allow direct error return from validator
func (m *Message) Validate(validator Validator) error {
	if m.bitmap == nil {
		return &MessageError{Message: "message not parsed, call Parse() first"}
	}

	if validator == nil {
		return nil
	}

	return validator.Validate(m)
}

// ValidateField validates a specific field.
// TODO: Implement once Spec is available.
func (m *Message) ValidateField(fieldNum int) error {
	// Field-level validation will be implemented with Spec support
	if !m.HasField(fieldNum) {
		return ErrFieldNotPresent
	}

	return nil
}

// parseAllFields eagerly parses all present fields and caches their cursors.
func (m *Message) parseAllFields() error {
	// Calculate starting offset after MTI and bitmap
	offset := 4 // MTI
	if m.bitmap.IsExtended() {
		offset += 16
	} else {
		offset += 8
	}

	// Parse each present field in order
	for _, fieldNum := range m.bitmap.PresentFields() {
		// Skip field 1 (bitmap itself)
		if fieldNum == 1 {
			continue
		}

		// Skip if field not in spec
		if _, ok := m.spec.Fields[fieldNum]; !ok {
			continue
		}

		// Parse field to get cursor
		cursor, err := m.parser.ParseField(m.buf, fieldNum, offset)
		if err != nil {
			return fmt.Errorf("failed to parse field %d: %w", fieldNum, err)
		}

		// Cache the cursor
		m.cursors[fieldNum] = cursor

		// Move to next field
		offset = cursor.NextOffset()
	}

	return nil
}

// isValidMTIStructure checks that MTI is 4 numeric ASCII digits.
func isValidMTIStructure(mti string) bool {
	if len(mti) != mtiLength {
		return false
	}

	for i := range mtiLength {
		if mti[i] < '0' || mti[i] > '9' {
			return false
		}
	}

	return true
}
