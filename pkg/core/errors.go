package core

import (
	"errors"
	"fmt"
)

// Predefined core errors.
var (
	ErrInvalidBitmap      = errors.New("invalid bitmap data")
	ErrFieldNotPresent    = errors.New("field not present")
	ErrInvalidFieldNumber = errors.New("invalid field number")
)

// MessageError wraps errors with additional context.
type MessageError struct {
	Message string
	Cause   error
}

func (e *MessageError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}

	return e.Message
}

func (e *MessageError) Unwrap() error {
	return e.Cause
}

// ErrMessageTooShort returns an error indicating the message is shorter than expected.
// expected: minimum required length
// actual: actual message length
func ErrMessageTooShort(expected, actual int) error {
	return &MessageError{
		Message: fmt.Sprintf("message too short: expected at least %d bytes, got %d", expected, actual),
	}
}

// ErrInvalidMTI returns an error for insufficient MTI length.
// length: actual MTI length
func ErrInvalidMTI(length int) error {
	return &MessageError{
		Message: fmt.Sprintf("invalid MTI: message must have at least 4 bytes for MTI, got %d", length),
	}
}

// ErrBitmapParseFailed returns an error for bitmap parsing failure.
// cause: underlying error
func ErrBitmapParseFailed(cause error) error {
	return &MessageError{
		Message: "failed to parse bitmap",
		Cause:   cause,
	}
}

// ErrInvalidMTIFormat returns an error for invalid MTI format (not 4 numeric digits).
// mti: MTI string
func ErrInvalidMTIFormat(mti string) error {
	return &MessageError{
		Message: fmt.Sprintf("invalid MTI format: MTI must be 4 numeric digits, got %q", mti),
	}
}

// ErrMissingRequiredField returns an error for a missing required field.
// fieldNum: ISO8583 field number.
func ErrMissingRequiredField(fieldNum int) error {
	return &MessageError{
		Message: fmt.Sprintf("missing required field: %d", fieldNum),
	}
}

// ErrInvalidFieldFormat returns an error for invalid field format.
func ErrInvalidFieldFormat(fieldNum int, reason string) error {
	return &MessageError{
		Message: fmt.Sprintf("invalid field %d format: %s", fieldNum, reason),
	}
}

// ErrInvalidPANChecksum returns an error for failed PAN Luhn checksum validation.
func ErrInvalidPANChecksum(fieldNum int) error {
	return &MessageError{
		Message: fmt.Sprintf("field %d: PAN failed Luhn checksum validation", fieldNum),
	}
}

// ErrInvalidFieldLength returns an error for invalid field length.
func ErrInvalidFieldLength(fieldNum, minLen, maxLen, actual int) error {
	return &MessageError{
		Message: fmt.Sprintf("field %d: length must be between %d and %d, got %d",
			fieldNum, minLen, maxLen, actual),
	}
}
