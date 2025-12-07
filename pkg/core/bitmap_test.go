package core_test

import (
	"testing"

	"github.com/hkumarmk/iso8583-lite/pkg/core"
)

func TestBitmap(t *testing.T) {
	t.Run("NewBitmap primary only", func(t *testing.T) {
		// Primary bitmap with field 2 set (0x40 = 01000000 in binary)
		data := []byte{0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

		bm, bytesRead, err := core.NewBitmap(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if bytesRead != 8 {
			t.Errorf("expected 8 bytes read, got %d", bytesRead)
		}

		if !bm.IsSet(2) {
			t.Error("expected field 2 to be set")
		}

		if bm.IsSet(3) {
			t.Error("expected field 3 to not be set")
		}
	})

	t.Run("NewBitmap with secondary", func(t *testing.T) {
		// Primary bitmap with field 1 set (indicating secondary bitmap)
		// 0x80 = 10000000 in binary (field 1 set)
		data := []byte{
			0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Primary
			0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Secondary (field 66 set)
		}

		bm, bytesRead, err := core.NewBitmap(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if bytesRead != 16 {
			t.Errorf("expected 16 bytes read, got %d", bytesRead)
		}

		if !bm.IsSet(1) {
			t.Error("expected field 1 to be set")
		}

		if !bm.IsSet(66) {
			t.Error("expected field 66 to be set")
		}
	})

	t.Run("NewBitmap invalid data", func(t *testing.T) {
		data := []byte{0x01, 0x02} // Too short

		_, _, err := core.NewBitmap(data)
		if err == nil {
			t.Error("expected error for invalid bitmap data")
		}
	})

	t.Run("Set and Unset", func(t *testing.T) {
		data := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		bm, _, _ := core.NewBitmap(data)

		// Set field 2
		bm.Set(2)

		if !bm.IsSet(2) {
			t.Error("expected field 2 to be set")
		}

		// Unset field 2
		bm.Unset(2)

		if bm.IsSet(2) {
			t.Error("expected field 2 to be unset")
		}

		// Set field 65 (should automatically set field 1)
		bm.Set(65)

		if !bm.IsSet(1) {
			t.Error("expected field 1 to be set when field 65 is set")
		}

		if !bm.IsSet(65) {
			t.Error("expected field 65 to be set")
		}
	})

	t.Run("PresentFields", func(t *testing.T) {
		data := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		bm, _, _ := core.NewBitmap(data)

		bm.Set(2)
		bm.Set(4)
		bm.Set(11)

		fields := bm.PresentFields()
		expectedFields := []int{2, 4, 11}

		if len(fields) != len(expectedFields) {
			t.Errorf("expected %d fields, got %d", len(expectedFields), len(fields))
		}

		for i, f := range expectedFields {
			if fields[i] != f {
				t.Errorf("expected field %d, got %d", f, fields[i])
			}
		}
	})

	t.Run("Bytes", func(t *testing.T) {
		data := []byte{0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		bm, _, _ := core.NewBitmap(data)

		bytes := bm.Bytes()
		if len(bytes) != 8 {
			t.Errorf("expected 8 bytes, got %d", len(bytes))
		}

		if bytes[0] != 0x40 {
			t.Errorf("expected first byte 0x40, got 0x%02x", bytes[0])
		}
	})
}
