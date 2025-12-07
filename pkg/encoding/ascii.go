package encoding

import (
	"errors"
	"fmt"
)

var (
	_ Encoder = (*asciiEncoder)(nil)

	//nolint:gochecknoglobals // ASCII is stateless and safe for concurrent use
	// ASCII is the default Encoder for ASCII encoding and decoding.
	ASCII Encoder = &asciiEncoder{}
)

const asciiMaxByte = 0x7F

var errNonASCIIByte = errors.New("non-ASCII byte")

// asciiEncoder implements Encoder for ASCII encoding.
type asciiEncoder struct{}

func (e *asciiEncoder) Encode(data []byte) ([]byte, error) {
	// ASCII encoding is a no-op for valid ASCII bytes
	for _, b := range data {
		if b > asciiMaxByte {
			return nil, fmt.Errorf("%w: 0x%X", errNonASCIIByte, b)
		}
	}

	out := make([]byte, len(data))
	copy(out, data)

	return out, nil
}

func (e *asciiEncoder) Decode(data []byte) ([]byte, int, error) {
	// ASCII decoding is a no-op for valid ASCII bytes
	for _, b := range data {
		if b > asciiMaxByte {
			return nil, 0, fmt.Errorf("%w: 0x%X", errNonASCIIByte, b)
		}
	}

	out := make([]byte, len(data))
	copy(out, data)

	return out, len(data), nil
}

func (e *asciiEncoder) Name() string {
	return "ASCII"
}
