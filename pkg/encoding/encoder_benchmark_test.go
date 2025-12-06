package encoding

import (
	"testing"
)

var (
	asciiTestData = []byte("1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

	// tlvTestData is a BER-TLV encoded byte slice:
	//   0x9F33 (Tag), 0x03 (Length), 0x01 0x02 0x03 (Value)
	//   0x95   (Tag), 0x02 (Length), 0xAA 0xBB (Value)
	tlvTestData    = []byte{0x9F, 0x33, 0x03, 0x01, 0x02, 0x03, 0x95, 0x02, 0xAA, 0xBB}
	ebcdicTestData = []byte("1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	binaryTestData = []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	bcdTestData    = []byte("1234567890")
	hexTestData    = []byte{0xDE, 0xAD, 0xBE, 0xEF, 0x01, 0x23, 0x45, 0x67}
)

func BenchmarkEBCDICEncode(b *testing.B) {
	enc := ebcdic037Encoder{}
	for i := 0; i < b.N; i++ {
		_, err := enc.Encode(ebcdicTestData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEBCDICDecode(b *testing.B) {
	enc := ebcdic037Encoder{}
	data, _ := enc.Encode(ebcdicTestData)
	for i := 0; i < b.N; i++ {
		_, _, err := enc.Decode(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBinaryEncode(b *testing.B) {
	enc := &binaryEncoder{}
	for i := 0; i < b.N; i++ {
		_, err := enc.Encode(binaryTestData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBinaryDecode(b *testing.B) {
	enc := &binaryEncoder{}
	data, _ := enc.Encode(binaryTestData)
	for i := 0; i < b.N; i++ {
		_, _, err := enc.Decode(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBCDEncode(b *testing.B) {
	enc := &bcdEncoder{}
	for i := 0; i < b.N; i++ {
		_, err := enc.Encode(bcdTestData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBCDDecode(b *testing.B) {
	enc := &bcdEncoder{}
	data, _ := enc.Encode(bcdTestData)
	for i := 0; i < b.N; i++ {
		_, _, err := enc.Decode(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHexEncode(b *testing.B) {
	enc := &hexEncoder{}
	for i := 0; i < b.N; i++ {
		_, err := enc.Encode(hexTestData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHexDecode(b *testing.B) {
	enc := &hexEncoder{}
	data, _ := enc.Encode(hexTestData)
	for i := 0; i < b.N; i++ {
		_, _, err := enc.Decode(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkASCIIEncode(b *testing.B) {
	enc := &asciiEncoder{}
	for i := 0; i < b.N; i++ {
		_, err := enc.Encode(asciiTestData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkASCIIDecode(b *testing.B) {
	enc := &asciiEncoder{}
	data, _ := enc.Encode(asciiTestData)
	for i := 0; i < b.N; i++ {
		_, _, err := enc.Decode(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTLVEncode(b *testing.B) {
	enc := &tlvEncoder{}
	for i := 0; i < b.N; i++ {
		_, err := enc.Encode(tlvTestData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTLVDecode(b *testing.B) {
	enc := &tlvEncoder{}
	data, _ := enc.Encode(tlvTestData)
	for i := 0; i < b.N; i++ {
		_, _, err := enc.Decode(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}
