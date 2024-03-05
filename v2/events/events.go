package events

import (
	"context"
	"io"
	"time"

	"github.com/alessandrofascini/log4go/v2/levels"
)

type Event struct {
	StartTime    time.Time
	CategoryName string
	Data         []any
	Level        levels.Level
	Context      context.Context
	Pid          int
}

type EventWriter interface {
	WriteEvent(e *Event) (n int, err error)
}

type EventWriterTo interface {
	WriteEventTo(w io.Writer, e *Event) (n int64, err error)
}
