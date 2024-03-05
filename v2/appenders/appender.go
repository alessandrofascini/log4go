package appenders

import (
	"io"

	"github.com/alessandrofascini/log4go/v2/events"
)

type Appender interface {
	io.Closer
	events.EventWriter
}
