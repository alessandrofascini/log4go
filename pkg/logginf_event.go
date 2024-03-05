package pkg

import "time"

type LoggingEvent struct {
	StartTime    time.Time
	CategoryName string
	Data         []any
	Level        Level
	Context      map[string]any
	Pid          int
}
