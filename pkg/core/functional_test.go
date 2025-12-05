package core_test

import (
	"encoding/hex"
	"testing"

	"github.com/hkumarmk/iso8583-lite/pkg/core"
)

func TestMessageFunctional(t *testing.T) {
	// Test Message structure:
	// MTI: 0200
	// Bitmap: fields 2, 3, 4, 11 present
	// Field 2: PAN "1234567890123456" (LLVAR)
	// Field 3: Processing Code "000000"
	// Field 4: Amount "000000001000"
	// Field 11: STAN "000001"

	msgHex := "303230307020000000000000313631323334353637383930313233343536303030303030303030303030303031303030303030303031"
	msgBytes, err := hex.DecodeString(msgHex)
	if err != nil {
		t.Fatalf("Failed to decode test message: %v", err)
	}

	// Parse the message using high-level Message API
	msg := core.NewMessage(msgBytes)

	// Parse MTI and bitmap
	err = msg.Parse()
	if err != nil {
		t.Fatalf("Failed to parse message: %v", err)
	}

	// Test Message.Bytes() returns original buffer
	if len(msg.Bytes()) != len(msgBytes) {
		t.Errorf("Expected message length %d, got %d", len(msgBytes), len(msg.Bytes()))
	}

	// Test MTI parsing
	expectedMTI := "0200"
	if msg.MTI().String() != expectedMTI {
		t.Errorf("Expected MTI %s, got %s", expectedMTI, msg.MTI().String())
	}

	// Test field presence via HasField (uses bitmap internally)
	if !msg.HasField(2) {
		t.Error("Expected field 2 to be present")
	}
	if !msg.HasField(3) {
		t.Error("Expected field 3 to be present")
	}
	if !msg.HasField(4) {
		t.Error("Expected field 4 to be present")
	}
	if !msg.HasField(11) {
		t.Error("Expected field 11 to be present")
	}
	if msg.HasField(5) {
		t.Error("Expected field 5 to NOT be present")
	}

	// Test PresentFields
	presentFields := msg.PresentFields()
	expectedFields := []int{0, 2, 3, 4, 11}
	if len(presentFields) != len(expectedFields) {
		t.Errorf("Expected %d present fields, got %d", len(expectedFields), len(presentFields))
	}
}
