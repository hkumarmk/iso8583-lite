package core

// Message represents a parsed ISO8583 message with zero-copy field access.
type Message struct {
	buf    []byte  // Raw message bytes
	bitmap *Bitmap // Bitmap indicating present fields (parsed lazily or on Parse)
	mti    string  // Message Type Indicator (field 0)
	fields map[int]Cursor
}

var _ MessageReader = (*Message)(nil)

// NewMessage creates a message wrapper without parsing.
// Use Parse() to parse MTI and bitmap.
func NewMessage(buf []byte) *Message {
	return &Message{
		buf:    buf,
		fields: make(map[int]Cursor),
	}
}

// Parse parses the MTI and bitmap with structural validation.
// Returns error for structural issues, nil on success.
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

	return nil
}

// Preload parses and caches all field cursors for hot access.
// Call after Parse() when you need high-performance repeated field access.
// Without this, fields are lazily accessed on-demand (fine for single reads).
// Returns self for method chaining.
// TODO: Implement once Spec is available for field parsing.
func (m *Message) Preload() *Message {
	// Field parsing will be implemented with Spec support
	// This will populate m.fields map with all present field cursors
	return m
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

	cursor, found := m.fields[fieldNum]
	if !found {
		return NewField(nil, false)
	}

	return NewField(cursor.Extract(m.buf), true)
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
