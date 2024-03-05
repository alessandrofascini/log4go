package pkg

import "time"

const (
	LocaleTime          = "15:04:05"
	ISO8601             = "2006-01-02T15:04:05.000"
	Iso8601WithTzOffset = "2006-01-02T15:04:05.000-07:00"
	AbsoluteTime        = "15:04:05.000"
	Datetime            = "2006/01/02-15:04:05"
)

func GetDateByFormat(format string) string {
	return time.Now().Format(format)
}
