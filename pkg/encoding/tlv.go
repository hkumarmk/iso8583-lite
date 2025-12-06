package encoding

import (
	"fmt"

	"github.com/euicc-go/bertlv"
)

// TODO: Consider replacing bertlv-based TLV encoder/decoder with a custom, zero-allocation implementation for high-performance systems.
// See docs/benchmark_results.md for details.
// TLVEncoder implements Encoder for BER-TLV encoding/decoding.
type tlvEncoder struct{}

var (
	_   Encoder = (*tlvEncoder)(nil)
	TLV Encoder = &tlvEncoder{}
)

// Encode encodes a slice of TLV objects into BER-TLV bytes.
// For ISO8583 field 55, input should be a raw TLV byte slice.

func (e *tlvEncoder) Encode(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return []byte{}, nil
	}
	var out []byte
	read := 0
	for read < len(data) {
		tlv := &bertlv.TLV{}
		err := tlv.UnmarshalBinary(data[read:])
		if err != nil {
			return nil, err
		}
		b, err := tlv.MarshalBinary()
		if err != nil {
			return nil, err
		}
		out = append(out, b...)
		read += len(b)
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
		tlv := &bertlv.TLV{}
		err := tlv.UnmarshalBinary(data[read:])
		if err != nil {
			return nil, read, fmt.Errorf("TLV decode error: %w", err)
		}
		b, err := tlv.MarshalBinary()
		if err != nil {
			return nil, read, fmt.Errorf("TLV marshal error: %w", err)
		}
		out = append(out, b...)
		read += len(b)
	}
	return out, read, nil
}

func (e *tlvEncoder) Name() string {
	return "TLV"
}
