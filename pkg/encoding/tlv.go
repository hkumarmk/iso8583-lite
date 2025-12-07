package encoding

import (
	"fmt"
)

// TODO: TLV implementation is currently optimized for flat TLV structures and short-form lengths (moov-io style).
// Limitations:
//   - Does not support constructed tags (nested TLVs)
//   - Does not support long-form lengths (multi-byte length fields)
//   - Limited error handling for malformed or deeply nested TLVs
//   - May not be fully compliant with all BER-TLV/EMV/ISO8583 extensions
// Future work:
//   - Extend parser for constructed tags and long-form lengths if required
//   - Improve error handling and validation
//   - Benchmark and validate against real-world data and full spec requirements
// Standards & References:
//   - ISO 8583: https://en.wikipedia.org/wiki/ISO_8583
//   - EMV Book 3: https://www.emvco.com/emv-technologies/specifications/
//   - ISO8583 Field 55: https://www.eftlab.com/knowledge-base/211-iso-8583-field-55-icc-system-related-data/
//   - BER-TLV: https://www.eftlab.com/knowledge-base/128-ber-tlv/
// See docs/benchmark_results.md for details.

// TLVEncoder implements Encoder for minimal BER-TLV encoding/decoding (flat TLV, short-form length).
type tlvEncoder struct{}

var (
	_ Encoder = (*tlvEncoder)(nil)

	//nolint:gochecknoglobals // TLV is stateless and safe for concurrent use
	// TLV is the default Encoder for BER-TLV encoding and decoding.
	TLV Encoder = &tlvEncoder{}
)

const (
	tlvTagMultiByteMask = 0x1F
	minimalTLVMinLen    = 2
)

// Encode encodes a slice of TLV objects into BER-TLV bytes (flat TLV, short-form length).
func (e *tlvEncoder) Encode(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return []byte{}, nil
	}

	var out []byte

	read := 0
	for read < len(data) {
		tag, _, value, next, err := ParseMinimalTLV(data[read:])
		if err != nil {
			return nil, err
		}

		out = append(out, EncodeMinimalTLV(tag, value)...)
		read += next
	}

	return out, nil
}

// Decode decodes BER-TLV bytes into raw bytes, returning the decoded data, number of bytes read, and error.
func (e *tlvEncoder) Decode(data []byte) ([]byte, int, error) {
	if len(data) == 0 {
		return []byte{}, 0, nil
	}

	var out []byte

	read := 0
	for read < len(data) {
		tag, _, value, next, err := ParseMinimalTLV(data[read:])
		if err != nil {
			return nil, read, fmt.Errorf("TLV decode error: %w", err)
		}

		out = append(out, EncodeMinimalTLV(tag, value)...)
		read += next
	}

	return out, read, nil
}

func (e *tlvEncoder) Name() string {
	return "TLV"
}

// MinimalTLV provides high-performance, minimal-allocation TLV parsing and encoding
// covering only flat TLV structures and short-form lengths (as in moov-io).
// See tlv.go for TODOs and limitations.

// ParseMinimalTLV parses a single TLV from data, returning tag, length, value, next offset, and error.
func ParseMinimalTLV(data []byte) ([]byte, int, []byte, int, error) {
	if len(data) < minimalTLVMinLen {
		return nil, 0, nil, 0, ErrTLVMalformed
	}
	// Tag parsing: single or multi-byte tag (only up to 2 bytes for minimal)
	tagLen := 1

	if data[0]&tlvTagMultiByteMask == tlvTagMultiByteMask {
		tagLen = 2
		if len(data) < tagLen+1 {
			return nil, 0, nil, 0, ErrTLVMalformed
		}
	}

	tag := data[:tagLen]
	// Length parsing: only short-form (1 byte)
	if len(data) < tagLen+1 {
		return nil, 0, nil, 0, ErrTLVMalformed
	}

	length := int(data[tagLen])
	// Value extraction
	if len(data) < tagLen+1+length {
		return nil, 0, nil, 0, ErrTLVMalformed
	}

	value := data[tagLen+1 : tagLen+1+length]
	next := tagLen + 1 + length

	return tag, length, value, next, nil
}

// EncodeMinimalTLV encodes tag and value into TLV bytes (short-form length only).
func EncodeMinimalTLV(tag, value []byte) []byte {
	out := make([]byte, len(tag)+1+len(value))
	copy(out, tag)
	out[len(tag)] = byte(len(value))
	copy(out[len(tag)+1:], value)

	return out
}

// ErrTLVMalformed is returned for malformed TLV input.
var ErrTLVMalformed = &TLVError{"malformed TLV"}

// TLVError represents an error related to TLV parsing or encoding.
type TLVError struct {
	msg string
}

func (e *TLVError) Error() string { return e.msg }
