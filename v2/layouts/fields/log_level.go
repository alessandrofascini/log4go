package fields

import "io"

type LogLevelField struct{}

func (l *LogLevelField) WriteTo(w io.Writer, reader *EventReader) (n int, err error) {
	return w.Write([]byte(reader.Level.LevelStr))
}
