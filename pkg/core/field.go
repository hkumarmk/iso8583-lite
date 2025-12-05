package core

import (
	"encoding/hex"
	"strconv"
)

// Field provides access to a single ISO8583 field with multiple type accessors.
type Field struct {
	data   []byte
	exists bool
}

var _ FieldAccessor = (*Field)(nil)

func NewField(data []byte, exists bool) *Field {
	return &Field{
		data:   data,
		exists: exists,
	}
}

func (f *Field) Exists() bool {
	return f.exists
}

func (f *Field) Bytes() []byte {
	if !f.exists {
		return nil
	}
	return f.data
}

func (f *Field) String() string {
	if !f.exists {
		return ""
	}
	return string(f.data)
}

func (f *Field) Int() int {
	val, _ := f.IntE()
	return val
}

func (f *Field) IntE() (int, error) {
	if !f.exists {
		return 0, ErrFieldNotPresent
	}
	return strconv.Atoi(f.String())
}

func (f *Field) Int64() int64 {
	val, _ := f.Int64E()
	return val
}

func (f *Field) Int64E() (int64, error) {
	if !f.exists {
		return 0, ErrFieldNotPresent
	}
	return strconv.ParseInt(f.String(), 10, 64)
}

func (f *Field) Hex() string {
	if !f.exists {
		return ""
	}
	return hex.EncodeToString(f.data)
}

func (f *Field) Len() int {
	return len(f.data)
}

// Deprecated: Use Len() instead for consistency with FieldAccessor interface
func (f *Field) Length() int {
	return f.Len()
}
