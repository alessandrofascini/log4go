package log4go

import (
	"os"
	"time"

	"github.com/alessandrofascini/log4go/pkg"

	loggerContext "github.com/alessandrofascini/log4go/internal/context"

	"github.com/alessandrofascini/log4go/internal/categories"

	"github.com/alessandrofascini/log4go/internal/appenders"
)

var ProcessPid = os.Getpid()

type Logger struct {
	loggerContext *loggerContext.LoggerContext
	name          string
	level         uint
	appenders     []*appenders.AppenderHandler
}

func NewLogger(c any) *Logger {
	switch category := c.(type) {
	case categories.Category:
		for _, handler := range category.Appenders {
			handler.Start()
		}
		return &Logger{
			loggerContext: loggerContext.NewLoggerContext(),
			name:          category.CategoryName,
			level:         category.LevelValue,
			appenders:     category.Appenders,
		}
	}
	panic("error while creating new logger")
}

// LoggerContext Wrapping

func (l *Logger) AddContext(key string, value any) {
	l.loggerContext.SetContext(key, value)
}

func (l *Logger) ChangeOneContext(key string, value any) {
	l.loggerContext.ChangeOneContext(key, value)
}

func (l *Logger) ChangeManyContext(c map[string]any) {
	for k, v := range c {
		l.ChangeOneContext(k, v)
	}
}

func (l *Logger) RemoveContext(key string) {
	l.loggerContext.RemoveContext(key)
}

func (l *Logger) ClearContext() {
	l.loggerContext.Clear()
}

// Log Methods

func (l *Logger) log(level pkg.Level, data []any) {
	event := &pkg.LoggingEvent{
		StartTime:    time.Now(),
		CategoryName: l.name,
		Data:         data,
		Level:        level,
		Context:      l.loggerContext.Consume(),
		Pid:          ProcessPid,
	}
	for _, handler := range l.appenders {
		handler.WriteOnChannel(event)
	}
}

func (l *Logger) Trace(args ...any) {
	if pkg.LevelTrace.Value < l.level {
		return
	}
	l.log(*pkg.LevelTrace, args)
}

func (l *Logger) Debug(args ...any) {
	if pkg.LevelTrace.Value < l.level {
		return
	}
	l.log(*pkg.LevelDebug, args)
}

func (l *Logger) Info(args ...any) {
	if pkg.LevelInfo.Value < l.level {
		return
	}
	l.log(*pkg.LevelInfo, args)
}

func (l *Logger) Warn(args ...any) {
	if pkg.LevelWarn.Value < l.level {
		return
	}
	l.log(*pkg.LevelWarn, args)
}

func (l *Logger) Error(args ...any) {
	if pkg.LevelError.Value < l.level {
		return
	}
	l.log(*pkg.LevelError, args)
}

func (l *Logger) Fatal(args ...any) {
	if pkg.LevelFatal.Value < l.level {
		return
	}
	l.log(*pkg.LevelFatal, args)
}

func (l *Logger) Terminate() {
	for _, handler := range l.appenders {
		handler.CloseChannel()
	}
}
