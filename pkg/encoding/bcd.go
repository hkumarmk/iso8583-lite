package encoding

import (
	"fmt"
)

var (
	_   Encoder = (*bcdEncoder)(nil)
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

// Encode encodes a digit string (ASCII bytes) into BCD bytes.
// Odd-length input is left-padded with '0'.
func (e *bcdEncoder) Encode(data []byte) ([]byte, error) {
	n := len(data)
	if n == 0 {
		return []byte{}, nil
	}
	if n%2 != 0 {
		data = append([]byte{'0'}, data...)
		n++
	}
	out := make([]byte, n/2)
	for i := 0; i < n; i += 2 {
		h := data[i]
		l := data[i+1]
		if h < '0' || h > '9' || l < '0' || l > '9' {
			return nil, fmt.Errorf("invalid digit in BCD input: %q%q", h, l)
		}
		out[i/2] = ((h - '0') << 4) | (l - '0')
	}
	return out, nil
}

// Decode decodes BCD bytes into a digit string (ASCII bytes).
func (e *bcdEncoder) Decode(data []byte) ([]byte, int, error) {
	if len(data) == 0 {
		return []byte{}, 0, nil
	}
	out := make([]byte, len(data)*2)
	for i, b := range data {
		h := (b >> 4) & 0x0F
		l := b & 0x0F
		if h > 9 || l > 9 {
			return nil, i, fmt.Errorf("invalid BCD digit: 0x%X", b)
		}
		out[i*2] = '0' + h
		out[i*2+1] = '0' + l
	}
	return out, len(data), nil
}

func (e *bcdEncoder) Name() string {
	return "BCD"
}
