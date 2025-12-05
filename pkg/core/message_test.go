package core

import (
	"strings"
	"testing"
)

func TestMessageParseErrors(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expectError string
	}{
		{
			name:        "empty message",
			data:        []byte{},
			expectError: "invalid MTI",
		},
		{
			name:        "too short for MTI",
			data:        []byte{0x30, 0x32},
			expectError: "invalid MTI: message must have at least 4 bytes for MTI, got 2",
		},
		{
			name:        "invalid MTI with non-numeric chars",
			data:        []byte{0x30, 0x32, 0x30, 0x41, 0x70, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // MTI: 020A
			expectError: "invalid MTI format: MTI must be 4 numeric digits",
		},
		{
			name:        "invalid MTI with special chars",
			data:        []byte{0x30, 0x32, 0x30, 0x2D, 0x70, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // MTI: 020-
			expectError: "invalid MTI format: MTI must be 4 numeric digits",
		},
		{
			name:        "too short for bitmap",
			data:        []byte{0x30, 0x32, 0x30, 0x30, 0x70, 0x20}, // MTI + 2 bitmap bytes
			expectError: "message too short: expected at least 12 bytes, got 6",
		},
		{
			name: "invalid bitmap",
			data: []byte{
				0x30, 0x32, 0x30, 0x30, // MTI
				0xFF, // Not enough bitmap data
			},
			expectError: "message too short",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := NewMessage(tt.data)
			err := msg.Parse()

			if err == nil {
				t.Fatal("Expected error but got nil")
			}

			if !strings.Contains(err.Error(), tt.expectError) {
				t.Errorf("Expected error containing %q, got %q", tt.expectError, err.Error())
			}

			t.Logf("Error message: %v", err)
		})
	}
}

func TestMessageParseSuccess(t *testing.T) {
	// Valid message with MTI and bitmap
	data := []byte{
		0x30, 0x32, 0x30, 0x30, // MTI: 0200
		0x70, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Primary bitmap
	}

	msg := NewMessage(data)
	err := msg.Parse()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if msg.MTI().String() != "0200" {
		t.Errorf("Expected MTI '0200', got '%s'", msg.MTI().String())
	}
}
