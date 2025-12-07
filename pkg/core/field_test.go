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

func TestFieldTreeStructure(t *testing.T) {
	// Create a parent field (simulating DE48 with composite data)
	parentData := []byte("ABCDEFGHIJKLMNOP")
	parent := NewField(parentData, true)

	// Create child fields (subfields within DE48)
	child1 := NewField(parentData[0:4], true)  // "ABCD"
	child2 := NewField(parentData[4:8], true)  // "EFGH"
	child3 := NewField(parentData[8:12], true) // "IJKL"

	// Set up parent-child relationships
	parent.SetSubfield(1, child1)
	parent.SetSubfield(2, child2)
	parent.SetSubfield(3, child3)

	// Test subfield access
	retrieved1 := parent.Subfield(1)
	if retrieved1.String() != "ABCD" {
		t.Errorf("Expected 'ABCD', got '%s'", retrieved1.String())
	}

	retrieved2 := parent.Subfield(2)
	if retrieved2.String() != "EFGH" {
		t.Errorf("Expected 'EFGH', got '%s'", retrieved2.String())
	}

	// Test non-existent subfield
	nonExistent := parent.Subfield(99)
	if nonExistent.Exists() {
		t.Error("Expected subfield 99 to not exist")
	}

	// Test HasSubfields
	if !parent.HasSubfields() {
		t.Error("Expected parent to have subfields")
	}

	emptyField := NewField(nil, false)
	if emptyField.HasSubfields() {
		t.Error("Expected empty field to have no subfields")
	}

	// Test multi-level hierarchy (3 levels)
	grandchild := NewField(child2.Bytes()[0:2], true) // "EF"
	child2.SetSubfield(1, grandchild)

	if grandchild.String() != "EF" {
		t.Errorf("Expected 'EF', got '%s'", grandchild.String())
	}

	// Access grandchild through parent
	retrieved := parent.Subfield(2).Subfield(1)
	if retrieved.String() != "EF" {
		t.Errorf("Expected 'EF', got '%s'", retrieved.String())
	}
}
