package pkg

import "strings"

const (
	All   = "ALL"
	Trace = "TRACE"
	Debug = "DEBUG"
	Info  = "INFO"
	Error = "ERROR"
	Warn  = "WARN"
	Fatal = "FATAL"
	Off   = "OFF"
)

type Level struct {
	Value    uint
	LevelStr string
	Color    Color
	Index    uint
}

var (
	LevelTrace = &Level{5000, Trace, *ColorBlue, 0}
	LevelDebug = &Level{10000, Debug, *ColorCyan, 1}
	LevelInfo  = &Level{20000, Info, *ColorGreen, 2}
	LevelWarn  = &Level{30000, Warn, *ColorYellow, 3}
	LevelError = &Level{40000, Error, *ColorRed, 4}
	LevelFatal = &Level{50000, Fatal, *ColorMagenta, 5}
)

func GetLevelByName(n string) Level {
	switch strings.ToUpper(n) {
	case Trace:
		return *LevelTrace
	case Debug:
		return *LevelDebug
	case Info:
		return *LevelInfo
	case Warn:
		return *LevelWarn
	case Error:
		return *LevelError
	case Fatal:
		return *LevelFatal
	}
	panic("cannot find a correct level")
}

func GetLevelValueByName(n string) uint {
	n = strings.ToUpper(n)
	switch n {
	case All:
		return 0
	case Off:
		return 60000
	default:
		return GetLevelByName(n).Value
	}
}
