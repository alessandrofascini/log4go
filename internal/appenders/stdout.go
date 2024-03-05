package appenders

import (
	"os"

	"github.com/alessandrofascini/log4go/pkg"
)

func StdoutAppender(_ pkg.AppenderConfig, layout *pkg.Layout) pkg.Appender {
	return func(event pkg.LoggingEvent) {
		_, err := os.Stdout.WriteString((*layout)(event))
		if err != nil {
			panic(err)
		}
	}
}
