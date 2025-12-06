package encoding

import "fmt"

var (
	_ Encoder = (*asciiEncoder)(nil)

	//nolint:gochecknoglobals // ASCII is stateless and safe for concurrent use
	ASCII Encoder = &asciiEncoder{}
)

// asciiEncoder implements Encoder for ASCII encoding.
type asciiEncoder struct{}

func (e *asciiEncoder) Encode(data []byte) ([]byte, error) {
	// ASCII encoding is a no-op for valid ASCII bytes
	for _, b := range data {
		if b > 0x7F {
			return nil, fmt.Errorf("non-ASCII byte: 0x%X", b)
		}
	}
	out := make([]byte, len(data))
	copy(out, data)

	return out, nil
}

func (e *asciiEncoder) Decode(data []byte) ([]byte, int, error) {
	// ASCII decoding is a no-op for valid ASCII bytes
	for _, b := range data {
		if b > 0x7F {
			return nil, 0, fmt.Errorf("non-ASCII byte: 0x%X", b)
		}
	}

	out := make([]byte, len(data))
	copy(out, data)

	return out, len(data), nil
}

func (e *asciiEncoder) Name() string {
	return "ASCII"
}
