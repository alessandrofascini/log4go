package fields

import (
	"io"
)

// Text Field

type TextField struct{ text []byte }

func NewTextField(b []byte) *TextField {
	f := &TextField{}
	f.text = make([]byte, len(b))
	copy(f.text, b)
	return f
}

func (t *TextField) WriteTo(w io.Writer, _ *EventReader) (n int, err error) {
	n, err = w.Write(t.text)
	return
}
