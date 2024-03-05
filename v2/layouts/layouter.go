package layouts

import (
	"github.com/alessandrofascini/log4go/v2/events"
)

// // Built-in layouts
// const (
// 	Basic              = "basic"
// 	Coloured           = "coloured"
// 	MessagePassThrough = "messagePassThrough"
// 	Dummy              = "dummy"
// 	Pattern            = "pattern"
// )

type Layouter interface {
	events.EventWriterTo
}

// func LayoutFactory(conf *Config) (Layouter, error) {
// 	switch conf.LayoutType {
// 	case Basic:
// 	case Coloured:
// 	case MessagePassThrough:
// 	case Dummy:
// 	case Pattern:
// 	}
// 	return nil, ErrUnkownLayout
// }
