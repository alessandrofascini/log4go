package layouts

import (
	"bytes"
	"io"
	"strings"

	"github.com/alessandrofascini/log4go/v2/events"
	"github.com/alessandrofascini/log4go/v2/layouts/fields"
)

// Layout define default layout options
// Implements Layouter interface
type Layout struct {
	// fields in pattern that may change at runtime (e.g. %m, %r). In golang std we call it verbs
	fields []fields.Field
}

// NewLayout: create new layout struct
// input: pattern as string
func NewLayout(p string) (*Layout, error) {
	reader := strings.NewReader(p)
	f := []fields.Field{}
	var buf bytes.Buffer
	for reader.Len() > 0 {
		r, _, err := reader.ReadRune()
		if err != nil {
			return nil, err
		}
		if r == '%' {
			if buf.Len() > 0 {
				f = append(f, fields.NewTextField(buf.Bytes()))
				buf.Reset()
			}
			if field, err := fields.NewField(reader); err != nil {
				return nil, err
			} else {
				f = append(f, field)
			}
		} else {
			buf.WriteRune(r)
		}
	}

	if buf.Len() > 0 {
		f = append(f, fields.NewTextField(buf.Bytes()))
	}

	l := &Layout{
		fields: make([]fields.Field, len(f)),
	}
	copy(l.fields, f)
	return l, nil
}

func (l *Layout) WriteEventTo(w io.Writer, e *events.Event) (n int64, err error) {
	var buf bytes.Buffer
	reader, err := fields.NewEventReader(e)
	if err != nil {
		return 0, err
	}
	for i := range l.fields {
		l.fields[i].WriteTo(&buf, reader)
	}
	m, err := w.Write(buf.Bytes())
	n = int64(m)
	return
}
