package fields

import (
	"bytes"
	"fmt"
	"io"
)

type padDirection = bool

type PaddingField struct {
	Field
	pad uint
	dir padDirection
}

func NewPaddingField(f Field, p int) *PaddingField {
	const (
		left  padDirection = true
		right padDirection = false
	)
	field := &PaddingField{f, uint(p), left}
	if p < 0 {
		field.dir = right
		field.pad = uint(p * -1)
	}
	return field
}

func (p *PaddingField) WriteTo(w io.Writer, reader *EventReader) (n int, err error) {
	var buf bytes.Buffer
	if _, err = p.Field.WriteTo(&buf, reader); err != nil {
		return
	}

	diff := buf.Len() - int(p.pad)
	fmt.Println(diff)
	if diff > 0 {
		m, err := buf.WriteTo(w)
		return int(m), err
	}

	diff *= -1
	b := make([]byte, p.pad)
	i := 0

	if p.dir {
		for i < diff {
			b[i] = ' '
			i++
		}
	}

	for buf.Len() > 0 {
		if b[i], err = buf.ReadByte(); err != nil {
			return 0, err
		}
		i++
	}

	if !p.dir {
		diff := uint(diff)
		for i := uint(0); i < diff; i++ {
			b[p.pad-i] = ' '
		}
	}

	return w.Write(b)
}
