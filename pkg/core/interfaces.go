package core

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

// MessageReader defines the complete interface for reading and validating ISO8583 messages.
type MessageReader interface {
	// Parse parses the MTI and bitmap from the message buffer.
	// Performs minimal Layer 1 (structural) validation.
	Parse() error

	// Preload parses and caches all field cursors for hot access.
	// Call after Parse() for high-performance repeated field access.
	Preload() *Message

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

// BitmapAccessor defines the interface for reading and modifying bitmap state.
type BitmapAccessor interface {
	// IsSet returns true if the specified field number is set in the bitmap.
	IsSet(fieldNum int) bool

	// PresentFields returns a slice of all field numbers present in the bitmap.
	PresentFields() []int

	// IsExtended returns true if the secondary bitmap is present.
	IsExtended() bool

	// Set marks the specified field as present.
	Set(fieldNum int)

	// Unset marks the specified field as absent.
	Unset(fieldNum int)

	// Bytes returns the bitmap as a byte slice for serialization.
	Bytes() []byte
}

// MessageBuilder defines the interface for constructing and modifying ISO8583 messages.
// This provides a fluent API for building messages with deferred validation.
type MessageBuilder interface {
	// SetMTI sets the Message Type Indicator.
	SetMTI(mti string) MessageBuilder

	// SetField sets a field value.
	SetField(fieldNum int, value interface{}) MessageBuilder

	// SetString sets a field from a string value.
	SetString(fieldNum int, value string) MessageBuilder

	// SetInt sets a field from an int value.
	SetInt(fieldNum int, value int) MessageBuilder

	// SetBytes sets a field from raw bytes.
	SetBytes(fieldNum int, value []byte) MessageBuilder

	// UnsetField removes a field.
	UnsetField(fieldNum int) MessageBuilder

	// Build finalizes the message and performs validation.
	// Returns the constructed message or error if validation fails.
	Build() (MessageReader, error)

	// BuildBytes finalizes the message and returns the serialized bytes.
	BuildBytes() ([]byte, error)
}
