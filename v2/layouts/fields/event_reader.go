package fields

import (
	"errors"

	"github.com/alessandrofascini/log4go/v2/events"
)

var (
	ErrNilPointer  = errors.New("nil pointer")
	ErrEmptyReader = errors.New("readed all data available")
)

type EventReader struct {
	*events.Event
	i int
}

func NewEventReader(e *events.Event) (*EventReader, error) {
	if e == nil {
		return nil, ErrNilPointer
	}
	return &EventReader{e, 0}, nil
}

func (r *EventReader) ReadData() (any, error) {
	if r.i == len(r.Data) {
		return nil, ErrEmptyReader
	}
	v := r.Data[r.i]
	r.i++
	return v, nil
}
