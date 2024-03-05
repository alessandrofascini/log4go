package appenders

import (
	"errors"

	"github.com/alessandrofascini/log4go/v2/events"
)

type Channeled struct {
	*Unchanneled
	chanSize uint
}

func NewChanneled(unchan *Unchanneled, chanSize uint) (*Channeled, error) {
	if chanSize == 0 {
		return nil, errors.New("chan size cannot be equal to zero")
	}
	return &Channeled{unchan, chanSize}, nil
}

func (u *Channeled) WriteEvent(e *events.Event) (int, error) {
	return 0, errors.New("implement me")
}

func (u *Channeled) Close() error {
	return errors.New("implement me")
}
