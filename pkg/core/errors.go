package core

import (
	"errors"
	"fmt"
)

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

func ErrMessageTooShort(expected, actual int) error {
	return &MessageError{
		Message: fmt.Sprintf("message too short: expected at least %d bytes, got %d", expected, actual),
	}
}

func ErrInvalidMTI(length int) error {
	return &MessageError{
		Message: fmt.Sprintf("invalid MTI: message must have at least 4 bytes for MTI, got %d", length),
	}
}

func ErrBitmapParseFailed(cause error) error {
	return &MessageError{
		Message: "failed to parse bitmap",
		Cause:   cause,
	}
}

func ErrInvalidMTIFormat(mti string) error {
	return &MessageError{
		Message: fmt.Sprintf("invalid MTI format: MTI must be 4 numeric digits, got %q", mti),
	}
}

func ErrMissingRequiredField(fieldNum int) error {
	return &MessageError{
		Message: fmt.Sprintf("missing required field: %d", fieldNum),
	}
}

func ErrInvalidFieldFormat(fieldNum int, reason string) error {
	return &MessageError{
		Message: fmt.Sprintf("invalid field %d format: %s", fieldNum, reason),
	}
}

func ErrInvalidPANChecksum(fieldNum int) error {
	return &MessageError{
		Message: fmt.Sprintf("field %d: PAN failed Luhn checksum validation", fieldNum),
	}
}

func ErrInvalidFieldLength(fieldNum, minLen, maxLen, actual int) error {
	return &MessageError{
		Message: fmt.Sprintf("field %d: length must be between %d and %d, got %d",
			fieldNum, minLen, maxLen, actual),
	}
}
