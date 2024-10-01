package appenders

import (
	"io"

	"github.com/alessandrofascini/log4go/v2/events"
)

type Unchanneled struct {
	writer io.WriteCloser
	layout events.EventWriterTo
}

func NewUnchanneled(w io.WriteCloser, l events.EventWriterTo) (*Unchanneled, error) {
	return &Unchanneled{w, l}, nil
}

func (u *Unchanneled) WriteEvent(e *events.Event) (int, error) {
	n, err := u.layout.WriteEventTo(u.writer, e)
	return int(n), err
}

func (u *Unchanneled) Close() error {
	return u.writer.Close()
}
