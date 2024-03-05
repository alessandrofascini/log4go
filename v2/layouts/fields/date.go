package fields

import (
	"io"
	"time"
)

const (
	LocaleTime          = "15:04:05"
	ISO8601             = "2006-01-02T15:04:05.000"
	ISO8601WithTzOffset = "2006-01-02T15:04:05.000-07:00"
	AbsoluteTime        = "15:04:05.000"
	Datetime            = "2006/01/02-15:04:05"
)

type DateField struct {
	layout string
}

func NewDateField(layout string) *DateField {
	switch layout {
	case "", "ISO8601":
		layout = ISO8601
	case "ISO8601_WITH_TZ_OFFSET":
		layout = ISO8601WithTzOffset
	case "ABSOLUTETIME":
		layout = AbsoluteTime
	case "DATETIME":
		layout = Datetime
	}
	return &DateField{layout}
}

func (d *DateField) WriteTo(w io.Writer, _ *EventReader) (n int, err error) {
	return w.Write([]byte(time.Now().Format(d.layout)))
}

type LocaleTimeField = DateField

func NewLocaleTimeField() *LocaleTimeField {
	return &LocaleTimeField{LocaleTime}
}
