package encoding

// BinaryEncoder implements Encoder for raw binary (no-op, pass-through).
type binaryEncoder struct{}

var (
	_      Encoder = (*binaryEncoder)(nil)
	Binary Encoder = &binaryEncoder{}
)

// Encode returns the input as-is (no encoding).
func (e *binaryEncoder) Encode(data []byte) ([]byte, error) {
	if data == nil {
		return []byte{}, nil
	}
	out := make([]byte, len(data))
	copy(out, data)
	return out, nil
}

// Decode returns the input as-is (no decoding).
func (e *binaryEncoder) Decode(data []byte) ([]byte, int, error) {
	if data == nil {
		return []byte{}, 0, nil
	}
	out := make([]byte, len(data))
	copy(out, data)
	return out, len(data), nil
}

func (e *binaryEncoder) Name() string {
	return "Binary"
}
