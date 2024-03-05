package fields

import (
	"bytes"
	"io"
)

type TruncationField struct {
	Field
	trunc int
}

func (t *TruncationField) WriteTo(w io.Writer, reader *EventReader) (n int, err error) {
	var buf bytes.Buffer
	if _, err = t.Field.WriteTo(&buf, reader); err != nil {
		return
	}
	b := buf.Bytes()
	m := len(b)
	start, end := 0, m
	if t.trunc > 0 {
		if t.trunc < m {
			end = t.trunc
		}
	} else {
		if t := m + t.trunc; t > 0 {
			start = t
		}
	}
	return w.Write(b[start:end])
}
