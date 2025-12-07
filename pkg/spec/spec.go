// Package spec defines the ISO8583 message specification, field types, encodings, and related structures.
package spec

// Spec defines the complete ISO8583 message specification.
// This is a singleton per message type - shared by all message instances.
type Spec struct {
	Name     string
	Version  string
	Defaults FieldDefaults
	Fields   map[int]*FieldSpec
}

// FieldDefaults defines default values for fields in a spec.
type FieldDefaults struct {
	Encoding EncodingType
	Padding  PaddingType
	PadChar  rune
}

// FieldSpec defines the specification for a single field.
type FieldSpec struct {
	Number      int
	Name        string
	Aliases     []string
	Type        FieldType
	Length      int // For fixed fields
	MaxLength   int // For variable fields
	DataType    DataType
	Encoding    EncodingType
	Padding     PaddingType
	PadChar     rune
	Description string
	Tag         string       // For TLV fields
	Children    []*FieldSpec // For composite fields (subfields)
}

// FieldType defines the type of field (fixed or variable length).
//
// Design Note: FieldType is implemented as an enum (int) with methods rather than
// as an interface with polymorphic implementations. This design choice prioritizes
// simplicity and performance for the fixed set of ISO 8583 field types.
//
// Rationale:
// - ISO 8583 has a well-defined, stable set of field types (Fixed, L, LL, LLL, Bitmap)
// - Enum approach is simpler: easier to serialize, compare, and reason about
// - Performance: no interface dispatch overhead, can be used as map keys
// - Parsing logic in Parser package using switch/case is straightforward
//
// Trade-off: Less extensible for custom field types, but this is acceptable since
// ISO 8583 field types are standardized and unlikely to change. If custom field
// types become necessary, we can refactor to an interface-based approach.
//
// See docs/decisions/adr-001-field-type-enum.md for full rationale.
type FieldType int

// FieldType enum values.
const (
	FieldTypeFixed  FieldType = iota // Fixed-length field
	FieldTypeL                       // Variable with 1-digit length indicator
	FieldTypeLL                      // Variable with 2-digit length indicator
	FieldTypeLLL                     // Variable with 3-digit length indicator
	FieldTypeBitmap                  // Bitmap field (special handling)
)

// String returns the string representation of FieldType.
func (ft FieldType) String() string {
	switch ft {
	case FieldTypeFixed:
		return "Fixed"
	case FieldTypeL:
		return "L"
	case FieldTypeLL:
		return "LL"
	case FieldTypeLLL:
		return "LLL"
	case FieldTypeBitmap:
		return "Bitmap"
	default:
		return "UnknownFieldType"
	}
}

// DataType defines the data type of field content.
type DataType int

// DataType enum values.
const (
	DataTypeNumeric DataType = iota
	DataTypeAlpha
	DataTypeAlphanumeric
	DataTypeAlphaNumericSpecial
	DataTypeBinary
)

// String returns the string representation of DataType.
func (dt DataType) String() string {
	switch dt {
	case DataTypeNumeric:
		return "Numeric"
	case DataTypeAlpha:
		return "Alpha"
	case DataTypeAlphanumeric:
		return "Alphanumeric"
	case DataTypeAlphaNumericSpecial:
		return "AlphaNumericSpecial"
	case DataTypeBinary:
		return "Binary"
	default:
		return "UnknownDataType"
	}
}

// EncodingType defines the encoding format for field data.
type EncodingType int

// EncodingType enum values.
const (
	EncodingASCII EncodingType = iota
	EncodingEBCDIC
	EncodingBCD
	EncodingBinary
)

// String returns the string representation of EncodingType.
func (et EncodingType) String() string {
	switch et {
	case EncodingASCII:
		return "ASCII"
	case EncodingEBCDIC:
		return "EBCDIC"
	case EncodingBCD:
		return "BCD"
	case EncodingBinary:
		return "Binary"
	default:
		return "UnknownEncoding"
	}
}

// PaddingType defines how fields should be padded.
type PaddingType int

// PaddingType enum values.
const (
	PaddingNone PaddingType = iota
	PaddingLeft
	PaddingRight
	PaddingCenter
)

// String returns the string representation of PaddingType.
func (pt PaddingType) String() string {
	switch pt {
	case PaddingNone:
		return "None"
	case PaddingLeft:
		return "Left"
	case PaddingRight:
		return "Right"
	case PaddingCenter:
		return "Center"
	default:
		return "UnknownPaddingType"
	}
}

// LengthIndicatorDigits returns the number of digits in the length indicator.
//
//nolint:exhaustive,mnd // We only care about L, LL, LLL here
func (ft FieldType) LengthIndicatorDigits() int {
	switch ft {
	case FieldTypeL:
		return 1
	case FieldTypeLL:
		return 2
	case FieldTypeLLL:
		return 3
	default:
		return 0
	}
}

// IsVariable returns true if the field type is variable length.
func (ft FieldType) IsVariable() bool {
	return ft >= FieldTypeL && ft <= FieldTypeLLL
}
