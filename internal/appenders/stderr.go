package appenders

import (
	"os"

	"github.com/alessandrofascini/log4go/pkg"
)

func StderrAppender(_ pkg.AppenderConfig, layout *pkg.Layout) pkg.Appender {
	return func(event pkg.LoggingEvent) {
		_, err := os.Stderr.WriteString((*layout)(event))
		if err != nil {
			panic(err)
		}
	}
}
