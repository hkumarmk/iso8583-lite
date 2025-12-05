package core

import (
	"testing"
)

func TestFieldAccessors(t *testing.T) {
	// Test with a simple field
	data := []byte("123456")
	field := NewField(data, true)

	// Test String
	if field.String() != "123456" {
		t.Errorf("Expected '123456', got '%s'", field.String())
	}

	// Test Int
	if field.Int() != 123456 {
		t.Errorf("Expected 123456, got %d", field.Int())
	}

	// Test IntE with valid data
	val, err := field.IntE()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if val != 123456 {
		t.Errorf("Expected 123456, got %d", val)
	}

	// Test with non-numeric data
	nonNumField := NewField([]byte("ABC123"), true)

	// Int should return 0 on error
	if nonNumField.Int() != 0 {
		t.Errorf("Expected 0 for invalid int, got %d", nonNumField.Int())
	}

	// IntE should return error
	_, err = nonNumField.IntE()
	if err == nil {
		t.Error("Expected error for invalid int")
	}

	// Test non-existent field
	emptyField := NewField(nil, false)
	if emptyField.Exists() {
		t.Error("Expected field to not exist")
	}
	if emptyField.String() != "" {
		t.Error("Expected empty string for non-existent field")
	}
	if emptyField.Int() != 0 {
		t.Error("Expected 0 for non-existent field")
	}

	_, err = emptyField.IntE()
	if err == nil {
		t.Error("Expected error for non-existent field")
	}
}
