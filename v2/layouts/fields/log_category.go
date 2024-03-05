package fields

import "io"

type LogCategoryField struct{}

func (l *LogCategoryField) WriteTo(w io.Writer, reader *EventReader) (n int, err error) {
	return w.Write([]byte(reader.CategoryName))
}
