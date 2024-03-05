package appenders

import (
	"fmt"
	"sync"

	"github.com/alessandrofascini/log4go/internal/layouts"
	"github.com/alessandrofascini/log4go/pkg"
)

// newAppenderExecutor This function create a single Appender
func newAppenderExecutor(config pkg.AppenderConfig) pkg.Appender {
	t := config.GetType()
	layout := layouts.NewLayout(config.GetLayout())
	switch t {
	case "stdout":
		return StdoutAppender(config, &layout)
	case "stderr":
		return StderrAppender(config, &layout)
	case "fileSync":
		return FileSyncAppender(config, &layout)
	case "file":
		return FileAppender(config, &layout)
	case "dateFile":
		return DateFileAppender(config, &layout)
	case "net":
		return NetworkAppender(config, &layout)
	case "http":
		return HttpAppender(config, &layout)
	}
	if externalAppendersPool[t] != nil {
		return (*externalAppendersPool[t])(&config)
	}
	panic(fmt.Sprintf("unsopported appender type: %q", t))
}

type AppenderHandler struct {
	CloseChannel   func()
	WriteOnChannel func(event *pkg.LoggingEvent)
	Start          func()
}

func NewAppenderHandler(config pkg.AppenderConfig, globalAppenderWaitGroup *sync.WaitGroup) *AppenderHandler {
	appender := newAppenderExecutor(config)
	channel := make(chan *pkg.LoggingEvent, config.GetChannelSize())
	isOpen := true
	isActive := false

	return &AppenderHandler{
		CloseChannel: func() {
			if isOpen {
				isOpen = false
				close(channel)
			}
		},
		WriteOnChannel: func(event *pkg.LoggingEvent) {
			if isOpen {
				channel <- event
			}
		},
		Start: func() {
			if isActive {
				return
			}
			isActive = true
			globalAppenderWaitGroup.Add(1)
			go func() {
				for event := range channel {
					appender(*event)
				}
				globalAppenderWaitGroup.Done()
			}()
		},
	}
}
