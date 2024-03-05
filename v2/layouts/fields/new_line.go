package fields

import (
	"io"
)

type singleCharField struct {
	b []byte
}

func newSingleCharFiled(b byte) *singleCharField {
	f := &singleCharField{make([]byte, 1)}
	f.b[0] = b
	return f
}

func (f *singleCharField) WriteTo(w io.Writer, _ *EventReader) (n int, err error) {
	return w.Write(f.b)
}

type NewLineField = singleCharField

func NewNewLineField() *NewLineField {
	return newSingleCharFiled('\n')
}

type PercentField = singleCharField

func NewPercentField() *singleCharField {
	return newSingleCharFiled('%')
}
