package encoding

// Encoder defines the interface for encoding and decoding field data.
type Encoder interface {
	Encode([]byte) ([]byte, error)
	Decode([]byte) ([]byte, int, error)
	Name() string
}
