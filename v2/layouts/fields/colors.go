package fields

import (
	"errors"
	"io"

	"github.com/alessandrofascini/log4go/v2/levels"
)

type Color struct {
	Name string
	Code string
}

var (
	ColorReset   = &Color{"reset", "\033[0m"}
	ColorRed     = &Color{"red", "\033[31m"}
	ColorGreen   = &Color{"green", "\033[32m"}
	ColorYellow  = &Color{"yellow", "\033[33m"}
	ColorBlue    = &Color{"blue", "\033[34m"}
	ColorMagenta = &Color{"magenta", "\033[35m"}
	ColorCyan    = &Color{"cyan", "\033[36m"}
	//ColorGray    = &Color{"gray", "\033[37m"}
	//ColorWhite   = &Color{"white", "\033[97m"}
)

type ColorField struct{}

func (c *ColorField) WriteTo(w io.Writer, reader *EventReader) (n int, err error) {
	var s string
	switch reader.Level.Value {
	case levels.LevelTrace.Value:
		s = ColorBlue.Code
	case levels.LevelDebug.Value:
		s = ColorCyan.Code
	case levels.LevelInfo.Value:
		s = ColorGreen.Code
	case levels.LevelWarn.Value:
		s = ColorYellow.Code
	case levels.LevelError.Value:
		s = ColorRed.Code
	case levels.LevelFatal.Value:
		s = ColorMagenta.Code
	default:
		return 0, errors.New("unkown level")
	}
	return w.Write([]byte(s))
}

type ColorResetField struct{}

func (c *ColorResetField) WriteTo(w io.Writer, _ *EventReader) (n int, err error) {
	return w.Write([]byte(ColorReset.Code))
}
