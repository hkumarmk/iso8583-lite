// Package core provides core ISO8583 message handling functionalities.
package core

import "encoding/binary"

// BitmapAccessor defines the interface for reading and modifying bitmap state.
type BitmapAccessor interface {
	// IsSet returns true if the specified field number is set in the bitmap.
	IsSet(fieldNum int) bool

	// PresentFields returns a slice of all field numbers present in the bitmap.
	PresentFields() []int

	// IsExtended returns true if the secondary bitmap is present.
	IsExtended() bool

	// Set marks the specified field as present.
	Set(fieldNum int)

	// Unset marks the specified field as absent.
	Unset(fieldNum int)

	// Bytes returns the bitmap as a byte slice for serialization.
	Bytes() []byte
}

// Bitmap represents the ISO8583 bitmap indicating which fields are present.
type Bitmap struct {
	primary   uint64
	secondary uint64
	extended  bool
}

var _ BitmapAccessor = (*Bitmap)(nil)

const (
	primaryBitmapLength     = 8
	secondaryBitmapLength   = 16
	primaryBitmapCapacity   = 64
	secondaryBitmapCapacity = 128
)

// NewBitmap parses the provided byte slice to construct a Bitmap instance according to the ISO8583 specification.
// It reads the primary bitmap (first 8 bytes) and, if the first bit is set, reads the secondary bitmap (next 8 bytes).
// Returns the Bitmap, the number of bytes read (8 or 16), and an error if the input data is invalid.
func NewBitmap(data []byte) (*Bitmap, int, error) {
	if len(data) < primaryBitmapLength {
		return nil, 0, ErrInvalidBitmap
	}

	bm := &Bitmap{
		primary: binary.BigEndian.Uint64(data[0:primaryBitmapLength]),
	}

	bytesRead := primaryBitmapLength

	// Field 1 set indicates secondary bitmap follows (per ISO8583 spec).
	// Secondary bitmap must be read even if all fields 65-128 are zero.
	if bm.IsSet(1) {
		if len(data) < secondaryBitmapLength {
			return nil, 0, ErrInvalidBitmap
		}

		bm.secondary = binary.BigEndian.Uint64(data[primaryBitmapLength:secondaryBitmapLength])
		bm.extended = true
		bytesRead = secondaryBitmapLength
	}

	return bm, bytesRead, nil
}

// IsSet returns true if the specified field number is set in the bitmap.
func (b *Bitmap) IsSet(fieldNum int) bool {
	if fieldNum < 1 || fieldNum > secondaryBitmapCapacity {
		return false
	}

	if fieldNum <= primaryBitmapCapacity {
		bit := uint64(1) << (primaryBitmapCapacity - fieldNum)

		return (b.primary & bit) != 0
	}

	if !b.extended {
		return false
	}

	bit := uint64(1) << (secondaryBitmapCapacity - fieldNum)

	return (b.secondary & bit) != 0
}

// Set marks the specified field as present in the bitmap.
func (b *Bitmap) Set(fieldNum int) {
	if fieldNum < 1 || fieldNum > secondaryBitmapCapacity {
		return
	}

	if fieldNum == 1 {
		b.extended = true
	}

	if fieldNum <= primaryBitmapCapacity {
		bit := uint64(1) << (primaryBitmapCapacity - fieldNum)
		b.primary |= bit
	} else {
		b.extended = true
		b.Set(1)

		bit := uint64(1) << (secondaryBitmapCapacity - fieldNum)
		b.secondary |= bit
	}
}

// Unset marks the specified field as absent in the bitmap.
func (b *Bitmap) Unset(fieldNum int) {
	if fieldNum < 1 || fieldNum > secondaryBitmapCapacity {
		return
	}

	if fieldNum <= primaryBitmapCapacity {
		bit := uint64(1) << (primaryBitmapCapacity - fieldNum)
		b.primary &^= bit
	} else {
		bit := uint64(1) << (secondaryBitmapCapacity - fieldNum)
		b.secondary &^= bit
	}
}

// Bytes returns the bitmap as a byte slice in big-endian order.
func (b *Bitmap) Bytes() []byte {
	if !b.extended {
		buf := make([]byte, primaryBitmapLength)
		binary.BigEndian.PutUint64(buf, b.primary)

		return buf
	}

	buf := make([]byte, secondaryBitmapLength)
	binary.BigEndian.PutUint64(buf[0:8], b.primary)
	binary.BigEndian.PutUint64(buf[8:16], b.secondary)

	return buf
}

// PresentFields returns a slice of integers representing the field numbers that are set in the bitmap.
func (b *Bitmap) PresentFields() []int {
	fields := make([]int, 0, primaryBitmapCapacity)

	for i := 1; i <= primaryBitmapCapacity; i++ {
		if b.IsSet(i) {
			fields = append(fields, i)
		}
	}

	if b.extended {
		for i := 65; i <= secondaryBitmapCapacity; i++ {
			if b.IsSet(i) {
				fields = append(fields, i)
			}
		}
	}

	return fields
}

// IsExtended returns true if the bitmap is in extended mode, indicating the presence of a secondary bitmap.
func (b *Bitmap) IsExtended() bool {
	return b.extended
}
