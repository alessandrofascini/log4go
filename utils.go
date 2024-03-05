package log4go

import (
	"fmt"
	"os"
	"sync"

	"github.com/alessandrofascini/log4go/pkg"

	configuration "github.com/alessandrofascini/log4go/internal/config"

	"github.com/alessandrofascini/log4go/internal/categories"

	"github.com/alessandrofascini/log4go/internal/appenders"
)

// useDefaultConfiguration this function configure Log4Go with default configuration
func useDefaultConfiguration() {
	Configure(map[string]any{
		"appenders":  map[string]any{"out": map[string]any{"type": "stdout"}},
		"categories": map[string]any{"default": map[string]any{"appenders": []string{"out"}, "level": "OFF"}},
	})
}

// validateLoggerConfig this function receive a map[string]any as configuration.
// if presents set global configurations
// return a more specific "implementation" of the map as LoggerConfig
func validateLoggerConfig(config map[string]any) pkg.LoggerConfig {
	const (
		AppendersField          = "appenders"
		CategoriesFieldName     = "categories"
		ConfigurationsFieldName = "configuration"
	)
	// Appenders
	switch v := config[AppendersField].(type) {
	case nil:
		panic(fmt.Sprintf("cannot find required field %q", AppendersField))
	default:
		panic(fmt.Sprintf("invalid type \"%T\" of field %q in logger config", v, AppendersField))
	case map[string]any:
		config[AppendersField] = validateAppendersConfig(v)
	}
	// Categories
	switch v := config[CategoriesFieldName].(type) {
	case nil:
		panic(fmt.Sprintf("cannot find required field %q", CategoriesFieldName))
	default:
		panic(fmt.Sprintf("invalid type \"%T\" of field %q in logger config", v, CategoriesFieldName))
	case map[string]any:
		config[CategoriesFieldName] = validateCategoriesConfig(v)
	}
	// Configurations
	switch settings := config[ConfigurationsFieldName].(type) {
	case map[string]any:
		for key, value := range settings {
			switch key {
			case "appenderChannelSize":
				if v, ok := anyToInt(value); ok {
					configuration.SetDefaultAppenderChannelSize(v)
				}
			case "internalAppenderChannelSize":
				if v, ok := anyToInt(value); ok {
					configuration.SetDefaultInternalChannelSize(v)
				}
			case "createFolderPerm":
				if v, ok := anyToInt(value); ok {
					configuration.SetCreateFolderPerm(os.FileMode(v))
				}
			}
		}
	}
	return config
}

func newCategoriesPool(config pkg.LoggerConfig, globalAppenderWaitGroup *sync.WaitGroup) *sync.Map {
	appenderHandlers := appenders.NewAppenderPool(config.GetAppenders(), globalAppenderWaitGroup)
	pool := sync.Map{}
	for name, categoryConfig := range config.GetCategories() {
		categoryConfig["name"] = name
		pool.Store(name, categories.NewCategory(categoryConfig, appenderHandlers))
	}
	return &pool
}

func anyToInt(value any) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	}
	return 0, false
}
