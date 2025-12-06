package core

import (
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

// Parse parses the MTI, bitmap, and all present fields (eager parsing).
// Returns an error if the message structure is invalid.
// For full validation including format and business rules, use Validate() with appropriate validators.
func (m *Message) Parse() error {
	if len(m.buf) < 4 {
		return ErrInvalidMTI(len(m.buf))
	}
	m.mti = string(m.buf[0:4])

	if !isValidMTIStructure(m.mti) {
		return ErrInvalidMTIFormat(m.mti)
	}

	if len(m.buf) < 12 {
		return ErrMessageTooShort(12, len(m.buf))
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
			return err
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
	if len(mti) != 4 {
		return false
	}
	for i := range 4 {
		if mti[i] < '0' || mti[i] > '9' {
			return false
		}
	}
	return true
}

func (m *Message) MTI() FieldAccessor {
	return NewField([]byte(m.mti), m.mti != "")
}

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

func (m *Message) HasField(fieldNum int) bool {
	if fieldNum == 0 {
		return true
	}
	if m.bitmap == nil {
		return false
	}
	return m.bitmap.IsSet(fieldNum)
}

func (m *Message) Bytes() []byte {
	return m.buf
}

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
// TODO: Implement once Spec is available
func (m *Message) ValidateField(fieldNum int) error {
	// Field-level validation will be implemented with Spec support
	if !m.HasField(fieldNum) {
		return ErrFieldNotPresent
	}
	return nil
}
