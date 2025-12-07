//nolint:mnd // Magic numbers are acceptable in encoding algorithms
package encoding

import (
	"errors"
	"fmt"
)

var (
	_ Encoder = (*bcdEncoder)(nil)

	//nolint:gochecknoglobals // BCD is stateless and safe for concurrent use
	// BCD is the Encoder for Binary Coded Decimal encoding and decoding.
	BCD Encoder = &bcdEncoder{}
)

// BCD (Binary Coded Decimal) encoding for ISO8583 fields.
//
// We implement BCD encoding/decoding in-house for clarity, performance, and minimalism.
// This avoids external dependencies and covers the standard ISO8583 use case:
//   - Packed BCD (2 digits per byte, left-aligned, pad with zero if odd)
//   - Only supports digit strings ('0'-'9') as per ISO8583 numeric field requirements
//   - Reference: ISO8583-1:2003, Section 7.2.4 (Numeric fields, packed BCD)
//   - See also: https://en.wikipedia.org/wiki/Binary-coded_decimal
//
// We do not support telecom TBCD, Excess-3, or other BCD variants, as they are not used in payment ISO8583 messages.
// If future requirements demand other variants, this can be extended or use a third-party library.

// bcdEncoder implements Encoder for BCD (Binary Coded Decimal).
type bcdEncoder struct{}

const digitsPerByte = 2

var errInvalidBCDDigit = errors.New("invalid BCD digit")

// Encode encodes a digit string (ASCII bytes) into BCD bytes.
// Odd-length input is left-padded with '0'.
func (e *bcdEncoder) Encode(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return []byte{}, nil
	}

	if len(data)%2 != 0 {
		paddedData := make([]byte, len(data)+1)
		paddedData[0] = '0'
		copy(paddedData[1:], data)
		data = paddedData
	}

	dataLength := len(data)

	out := make([]byte, dataLength/digitsPerByte)
	for idx := 0; idx < dataLength; idx += digitsPerByte {
		high := data[idx]
		low := data[idx+1]

		if high < '0' || high > '9' || low < '0' || low > '9' {
			return nil, fmt.Errorf("%w: %q%q", errInvalidBCDDigit, high, low)
		}

		out[idx/2] = ((high - '0') << 4) | (low - '0')
	}

	return out, nil
}

// Decode decodes BCD bytes into a digit string (ASCII bytes).
func (e *bcdEncoder) Decode(data []byte) ([]byte, int, error) {
	if len(data) == 0 {
		return []byte{}, 0, nil
	}

	out := make([]byte, len(data)*digitsPerByte)

	for i, b := range data {
		h := (b >> 4) & 0x0F
		l := b & 0x0F

		if h > 9 || l > 9 {
			return nil, i, fmt.Errorf("%w: 0x%X", errInvalidBCDDigit, b)
		}

		out[i*digitsPerByte] = '0' + h
		out[i*digitsPerByte+1] = '0' + l
	}

	return out, len(data), nil
}

func (e *bcdEncoder) Name() string {
	return "BCD"
}
