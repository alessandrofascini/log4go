package v2

import (
	"errors"
	"github.com/alessandrofascini/log4go/v2/events"
)

type Logger struct {
	appenders []events.EventWriter
}

func NewLogger(*Config) (*Logger, error) {
	return nil, errors.New("NewLogger: implement me")
}
