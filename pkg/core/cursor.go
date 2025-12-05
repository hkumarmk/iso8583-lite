package core

// Cursor represents a zero-copy position in the message buffer.
type Cursor struct {
	Start int
	End   int
}

func (c Cursor) Length() int {
	return c.End - c.Start
}

func (c Cursor) IsValid() bool {
	return c.Start >= 0 && c.End >= c.Start
}

func (c Cursor) IsEmpty() bool {
	return c.Start == c.End
}

func (c Cursor) Extract(buf []byte) []byte {
	if !c.IsValid() || c.End > len(buf) {
		return nil
	}
	return buf[c.Start:c.End]
}
