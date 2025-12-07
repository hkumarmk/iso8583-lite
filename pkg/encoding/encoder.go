// Package encoding provides interfaces and implementations for encoding and decoding byte slices,
// commonly used for data transformation in ISO8583 message processing. The Encoder interface defines
// methods for encoding and decoding byte slices, as well as retrieving the encoder's name.
package encoding

// Encoder defines an interface for encoding and decoding byte slices.
type Encoder interface {
	// Encode encodes the provided data and returns the encoded result or an error.
	Encode(data []byte) ([]byte, error)
	// Decode decodes the provided data and returns the decoded result, the number of bytes consumed, and an error.
	Decode(data []byte) ([]byte, int, error)
	// Name returns the name of the encoder implementation.
	Name() string
}
