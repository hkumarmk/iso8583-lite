package core

import "fmt"

// Validator defines the interface for ISO8583 message validation.
// Implementations can validate at different layers:
// - Layer 1: Structural validation (field parsing, boundaries)
// - Layer 1.5: Format validation (data types, patterns, mandatory fields)
// - Layer 2: Business validation (Luhn checks, expiration, limits).
type Validator interface {
	// Validate checks the message and returns error if validation fails.
	// The error should be descriptive and indicate which field/rule failed.
	Validate(msg MessageReader) error
}

// MessageBuilder defines the interface for constructing and modifying ISO8583 messages.
// This provides a fluent API for building messages with deferred validation.
type MessageBuilder interface {
	// SetMTI sets the Message Type Indicator.
	SetMTI(mti string) MessageBuilder

	// SetField sets a field value.
	SetField(fieldNum int, value any) MessageBuilder

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

// ValidatorFunc is a function adapter for the Validator interface.
type ValidatorFunc func(MessageReader) error

// Validate calls the underlying function.
func (f ValidatorFunc) Validate(msg MessageReader) error {
	return f(msg)
}

// CompositeValidator chains multiple validators in sequence.
// Stops at first validation error.
type CompositeValidator struct {
	validators []Validator
}

// NewCompositeValidator creates a composite validator from multiple validators.
func NewCompositeValidator(validators ...Validator) *CompositeValidator {
	return &CompositeValidator{
		validators: validators,
	}
}

// Validate runs all validators in sequence.
//
//nolint:wrapcheck // Allow direct error return from validators
func (c *CompositeValidator) Validate(msg MessageReader) error {
	for _, v := range c.validators {
		if err := v.Validate(msg); err != nil {
			return err
		}
	}

	return nil
}

// StructuralValidator validates message structure (Layer 1).
// - All fields can be parsed according to spec
// - Field boundaries are valid
// - Variable length indicators are correct.
type StructuralValidator struct {
	// Spec defines field formats for structural validation
	// TODO: Will be implemented when Spec is available
}

// NewStructuralValidator creates a new structural validator.
func NewStructuralValidator() *StructuralValidator {
	return &StructuralValidator{}
}

// Validate performs structural validation.
func (v *StructuralValidator) Validate(_ MessageReader) error {
	// TODO: Implement structural validation
	// This will parse all fields and catch:
	// - Truncated fields
	// - Invalid length indicators
	// - Field boundary violations
	return nil
}

// FormatValidator validates field formats (Layer 1.5).
// - Mandatory fields present
// - Data types correct (numeric/alpha/alphanumeric)
// - Length constraints satisfied
// - Date/time patterns valid.
type FormatValidator struct {
	// TODO: Add spec-based validation rules
}

// NewFormatValidator creates a new format validator.
func NewFormatValidator() *FormatValidator {
	return &FormatValidator{}
}

// Validate performs format validation.
func (v *FormatValidator) Validate(_ MessageReader) error {
	// TODO: Implement format validation
	// - Check mandatory fields
	// - Validate data types per spec
	// - Check length constraints
	// - Validate patterns (dates, times, amounts)
	return nil
}

// BusinessValidator validates business rules (Layer 2).
type BusinessValidator struct {
	rules []ValidationRule
}

// ValidationRule defines a business rule validation.
type ValidationRule interface {
	// Check validates the rule and returns error if it fails.
	Check(msg MessageReader) error
}

// NewBusinessValidator creates a business validator with given rules.
func NewBusinessValidator(rules ...ValidationRule) *BusinessValidator {
	return &BusinessValidator{
		rules: rules,
	}
}

// Validate performs business rule validations.
func (v *BusinessValidator) Validate(msg MessageReader) error {
	for _, rule := range v.rules {
		if err := rule.Check(msg); err != nil {
			return fmt.Errorf("business rule validation failed: %w", err)
		}
	}

	return nil
}

// Common validation rules

// RequiredFieldsRule validates that required fields are present.
type RequiredFieldsRule struct {
	fields []int
}

// NewRequiredFieldsRule creates a rule for required fields.
func NewRequiredFieldsRule(fields ...int) *RequiredFieldsRule {
	return &RequiredFieldsRule{fields: fields}
}

// Check validates that all required fields are present.
func (r *RequiredFieldsRule) Check(msg MessageReader) error {
	for _, fieldNum := range r.fields {
		if !msg.HasField(fieldNum) {
			return ErrMissingRequiredField(fieldNum)
		}
	}

	return nil
}

// NumericFieldRule validates that fields contain only numeric characters.
type NumericFieldRule struct {
	fields []int
}

// NewNumericFieldRule creates a rule for numeric fields.
func NewNumericFieldRule(fields ...int) *NumericFieldRule {
	return &NumericFieldRule{fields: fields}
}

// Check validates that specified fields are numeric.
func (r *NumericFieldRule) Check(msg MessageReader) error {
	for _, fieldNum := range r.fields {
		if !msg.HasField(fieldNum) {
			continue // Skip if field not present
		}

		data := msg.Field(fieldNum).Bytes()
		for _, b := range data {
			if b < '0' || b > '9' {
				return ErrInvalidFieldFormat(fieldNum, "must be numeric")
			}
		}
	}

	return nil
}

// LuhnCheckRule validates PAN using Luhn algorithm.
type LuhnCheckRule struct {
	fieldNum int
}

// NewLuhnCheckRule creates a rule for Luhn check validation on a field.
func NewLuhnCheckRule(fieldNum int) *LuhnCheckRule {
	return &LuhnCheckRule{fieldNum: fieldNum}
}

// Check performs the Luhn check validation.
func (r *LuhnCheckRule) Check(msg MessageReader) error {
	if !msg.HasField(r.fieldNum) {
		return nil // Skip if field not present
	}

	pan := msg.Field(r.fieldNum).String()
	if !luhnCheck(pan) {
		return ErrInvalidPANChecksum(r.fieldNum)
	}

	return nil
}

const (
	luhnParityDivisor    = 2
	luhnDoubleSubtractor = 9
)

// luhnCheck validates a number using the Luhn algorithm.
func luhnCheck(number string) bool {
	var sum int

	parity := len(number) % luhnParityDivisor

	for i, digit := range number {
		if digit < '0' || digit > '9' {
			return false
		}

		d := int(digit - '0')
		if i%luhnParityDivisor == parity {
			d *= luhnParityDivisor
			if d > luhnDoubleSubtractor {
				d -= luhnDoubleSubtractor
			}
		}

		sum += d
	}

	return sum%10 == 0
}

// FieldLengthRule validates field length constraints.
type FieldLengthRule struct {
	fieldNum int
	minLen   int
	maxLen   int
}

// NewFieldLengthRule creates a rule for field length constraints.
func NewFieldLengthRule(fieldNum, minLen, maxLen int) *FieldLengthRule {
	return &FieldLengthRule{
		fieldNum: fieldNum,
		minLen:   minLen,
		maxLen:   maxLen,
	}
}

// Check validates the field length constraints.
func (r *FieldLengthRule) Check(msg MessageReader) error {
	if !msg.HasField(r.fieldNum) {
		return nil
	}

	length := len(msg.Field(r.fieldNum).Bytes())
	if length < r.minLen || length > r.maxLen {
		return ErrInvalidFieldLength(r.fieldNum, r.minLen, r.maxLen, length)
	}

	return nil
}
