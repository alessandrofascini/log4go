package levels

import (
	"errors"
	"strings"
)

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
	Value    int
	LevelStr string
	Index    uint
}

var (
	LevelTrace = &Level{5000, Trace, 0}
	LevelDebug = &Level{10000, Debug, 1}
	LevelInfo  = &Level{20000, Info, 2}
	LevelWarn  = &Level{30000, Warn, 3}
	LevelError = &Level{40000, Error, 4}
	LevelFatal = &Level{50000, Fatal, 5}
)

// Can I do this better?

func GetLevelByName(levelName string) (*Level, error) {
	switch strings.ToUpper(levelName) {
	case Trace:
		return LevelTrace, nil
	case Debug:
		return LevelDebug, nil
	case Info:
		return LevelInfo, nil
	case Warn:
		return LevelWarn, nil
	case Error:
		return LevelError, nil
	case Fatal:
		return LevelFatal, nil
	}
	return nil, errors.New("cannot find level with this name: " + levelName)
}

func GetLevelValueByName(n string) int {
	n = strings.ToUpper(n)
	switch n {
	case All:
		return 0
	case Off:
		return 60000
	}
	level, err := GetLevelByName(n)
	if err != nil {
		return -1
	}
	return level.Value
}
