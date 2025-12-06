package parser

// Cursor represents a zero-copy position in the message buffer.
// It tracks the start and end positions of a field's data within the raw message bytes.
type Cursor struct {
	Start int
	End   int
}

// Length returns the length of the data segment.
func (c Cursor) Length() int {
	return c.End - c.Start
}

// Extract returns the data segment from the buffer.
// Returns nil if the cursor is invalid or out of bounds.
func (c Cursor) Extract(buf []byte) []byte {
	if c.Start < 0 || c.End > len(buf) || c.Start >= c.End {
		return nil
	}
	return buf[c.Start:c.End]
}

// NextOffset returns the offset where the next field should start.
// This is simply the End position of the current cursor.
func (c Cursor) NextOffset() int {
	return c.End
}
