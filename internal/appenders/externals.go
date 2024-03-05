package appenders

import "github.com/alessandrofascini/log4go/pkg"

// externalAppendersPool : this variable is used to store custom appenders. when logger is configured, this values
var externalAppendersPool = make(map[string]*func(config *pkg.AppenderConfig) pkg.Appender)

func AddExternalAppender(name string, f *func(config *pkg.AppenderConfig) pkg.Appender) {
	externalAppendersPool[name] = f
}

func ClearAppenderPool() {
	externalAppendersPool = make(map[string]*func(config *pkg.AppenderConfig) pkg.Appender)
}
