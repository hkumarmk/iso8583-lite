package core_test

import (
	"testing"

	"github.com/hkumarmk/iso8583-lite/pkg/core"
)

func TestCursor(t *testing.T) {
	t.Run("Length", func(t *testing.T) {
		c := core.Cursor{Start: 10, End: 20}
		if c.Length() != 10 {
			t.Errorf("expected length 10, got %d", c.Length())
		}
	})

	t.Run("IsValid", func(t *testing.T) {
		tests := []struct {
			name   string
			cursor core.Cursor
			want   bool
		}{
			{"valid", core.Cursor{Start: 0, End: 10}, true},
			{"valid empty", core.Cursor{Start: 5, End: 5}, true},
			{"invalid negative start", core.Cursor{Start: -1, End: 10}, false},
			{"invalid end before start", core.Cursor{Start: 10, End: 5}, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.cursor.IsValid(); got != tt.want {
					t.Errorf("IsValid() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("IsEmpty", func(t *testing.T) {
		c1 := core.Cursor{Start: 5, End: 5}
		if !c1.IsEmpty() {
			t.Error("expected cursor to be empty")
		}

		c2 := core.Cursor{Start: 5, End: 10}
		if c2.IsEmpty() {
			t.Error("expected cursor to not be empty")
		}
	})

	t.Run("Extract", func(t *testing.T) {
		buf := []byte("Hello, World!")
		c := core.Cursor{Start: 7, End: 12}

		extracted := c.Extract(buf)
		expected := "World"
		if string(extracted) != expected {
			t.Errorf("expected %q, got %q", expected, string(extracted))
		}
	})

	t.Run("Extract out of bounds", func(t *testing.T) {
		buf := []byte("Hello")
		c := core.Cursor{Start: 0, End: 100}

		extracted := c.Extract(buf)
		if extracted != nil {
			t.Error("expected nil for out of bounds cursor")
		}
	})
}
