package fields

import (
	"io"
	"strconv"
)

type ProcessIdField struct{}

func (p *ProcessIdField) WriteTo(w io.Writer, reader *EventReader) (n int, err error) {
	return w.Write([]byte(strconv.Itoa(reader.Pid)))
}
