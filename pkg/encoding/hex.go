package encoding

import (
	"encoding/hex"
)

// HexEncoder implements Encoder for hexadecimal string encoding.
type hexEncoder struct{}

var (
	_   Encoder = (*hexEncoder)(nil)
	Hex Encoder = &hexEncoder{}
)

// Encode encodes bytes as a lowercase hex string (ASCII bytes).
func (e *hexEncoder) Encode(data []byte) ([]byte, error) {
	if data == nil {
		return []byte{}, nil
	}
	out := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(out, data)
	return out, nil
}

// Decode decodes a hex string (ASCII bytes) to raw bytes.
func (e *hexEncoder) Decode(data []byte) ([]byte, int, error) {
	if data == nil {
		return []byte{}, 0, nil
	}
	out := make([]byte, hex.DecodedLen(len(data)))
	n, err := hex.Decode(out, data)
	if err != nil {
		return nil, 0, err
	}
	// n is the number of bytes written to out, but we want to return the number of bytes read from input (len(data))
	return out[:n], len(data), nil
}

func (e *hexEncoder) Name() string {
	return "Hex"
}
