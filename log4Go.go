// Package log4go
//
// This is a golang Logger inspired by log4js
package log4go

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/alessandrofascini/log4go/internal/appenders"
	"github.com/alessandrofascini/log4go/internal/layouts"
	"github.com/alessandrofascini/log4go/pkg"
)

var isConfigured = false

var categoriesPool = &sync.Map{}
var loggerPool = &sync.Map{}

var appendersWaitGroup sync.WaitGroup

// AddLayout Use this function to creating a new custom layouts.
// If you use a predefined key (e.g. basic, pattern and others), Log4Go will not use the custom layout.
// If you call this function twice, with the same key but different functions,
// Log4Go will use the last provided definition
func AddLayout(key string, f *func(*pkg.LayoutConfig) pkg.Layout) {
	if !isConfigured {
		layouts.AddExternalLayout(key, f)
	}
}

// AddAppender Use this function to making new custom appenders.
// You can't use a default appender name like dateFile, file and so on.
// If you call this function twice with the same key, the last function provided will be retained
func AddAppender(key string, f *func(config *pkg.AppenderConfig) pkg.Appender) {
	if !isConfigured {
		appenders.AddExternalAppender(key, f)
	}
}

// IsConfigured Is your Log4Go configured?
func IsConfigured() bool {
	return isConfigured
}

// Configure This function allow you to configure your logger.
func Configure(c any) error {
	if isConfigured {
		panic(fmt.Errorf("you've already configured your logger. this action can be executed only one time"))
	}
	configMap := make(map[string]any)

	// check if is a string or a map of string
	switch t := c.(type) {
	case string: // json file
		var err error
		file, err := os.ReadFile(t)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(file, &configMap); err != nil {
			return err
		}
	case map[string]any:
		// ok
		configMap = t
	default:
		return fmt.Errorf("invalid configure type")
	}

	// validate configuration passed
	config := validateLoggerConfig(configMap)

	// create categoriesPool
	categoriesPool = newCategoriesPool(config, &appendersWaitGroup)

	// now logger is configured
	isConfigured = true
	return nil
}

// GetLogger this function returns a Log4GoLogger.
// its get an array of strings, but it used only the first
// if you don't pass anything to this method, it returns a Log4GoLogger with name "defaults"
func GetLogger(names ...string) *Logger {
	if !isConfigured {
		useDefaultConfiguration()
	}
	var loggerName string
	if len(names) == 0 {
		loggerName = "default"
	} else {
		loggerName = names[0]
	}
	logger, ok := loggerPool.Load(loggerName)
	if !ok {
		// Create new logger from a category and then remove category from pool
		category, ok := categoriesPool.LoadAndDelete(loggerName)
		if !ok {
			panic(fmt.Errorf("%q logger does not exists", loggerName))
		}
		newLogger := NewLogger(category)
		loggerPool.Store(loggerName, newLogger)
		return newLogger
	}
	switch t := logger.(type) {
	case *Logger:
		return t
	}
	panic(fmt.Errorf("Log4Go internal error: %q is not a logger", loggerName))
}

// Shutdown This function must be pushed at the end of your main function.
// With this function you will be sure that all your logs were write on their files
func Shutdown() {
	loggerPool.Range(func(_, value any) bool {
		switch logger := value.(type) {
		case pkg.ILogger:
			logger.Terminate()
		}
		return true
	})
	appendersWaitGroup.Wait()
	deconfigure()
}

func deconfigure() {
	isConfigured = false
	categoriesPool = new(sync.Map)
	loggerPool = new(sync.Map)
	appendersWaitGroup = sync.WaitGroup{}
	// empty space
	appenders.ClearAppenderPool()
	layouts.ClearLayoutPool()
}
