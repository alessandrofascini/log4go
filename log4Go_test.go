package log4go

import (
	"fmt"
	"testing"

	"github.com/alessandrofascini/log4go/pkg"
)

func TestOrchestrator_Configure(t *testing.T) {
	Configure(map[string]any{
		"appenders": map[string]any{
			"out": map[string]any{
				"type": "stdout",
				"layout": map[string]any{
					"type": "coloured",
				},
			},
		},
		"categories": map[string]any{
			"cheese": map[string]any{
				"appenders": []string{
					"out",
				},
				"level": "error",
			},
			"default": map[string]any{
				"appenders": []string{
					"out",
				},
				"level": "all",
			},
		},
	})
	l := GetLogger("cheese")
	l.Info("info")
	l.Fatal("a fatal error")

	logger := GetLogger()
	logger.Info("info")
	logger.Fatal("a fatal error")
	logger.Info("5")
	logger.Info("4")
	logger.Info("3")
	logger.Info("2")
	logger.Info("1")
	logger.Info("hot file")
	for i := 0; i < 10; i++ {
		logger.Trace("TRACE log", i)
		logger.Debug("DEBUG log", i)
		logger.Info("INFO Log", i)
		logger.Warn("WARN Log", i)
		logger.Error("ERROR Log", i)
		logger.Fatal("FATAL log", i)
	}
	logger.Info("Before Shutdown")
	l.Fatal("prima della fine")
	Shutdown()
	logger.Info("After Shutdown")
}

func TestAddLayout(t *testing.T) {
	customLayout := func(config *pkg.LayoutConfig) pkg.Layout {
		return func(event pkg.LoggingEvent) string {
			return fmt.Sprintf("{ %q : %q, %q : %q }\n", "timestamp", event.StartTime, "data", event.Data)
		}
	}
	AddLayout("json", &customLayout)
	Configure(map[string]any{
		"appenders": map[string]any{
			"out": map[string]any{
				"type": "stdout",
				"layout": map[string]any{
					"type": "json",
				},
			},
		},
		"categories": map[string]any{
			"default": map[string]any{
				"appenders": []string{
					"out",
				},
				"level": "all",
			},
		},
	})
	defer func() {
		if recover() != nil {
			t.FailNow()
		}
	}()
	Shutdown()
}

func TestAddAppender(t *testing.T) {
	// t.SkipNow()
	customAppender := func(config *pkg.AppenderConfig) pkg.Appender {
		return nil
	}
	t.Log(&customAppender)
	AddAppender("test", &customAppender)

	// Shutdown()
}

func BenchmarkLogging(b *testing.B) {
	// Configure(map[string]any{
	// 	"appenders": map[string]any{
	// 		"out": map[string]any{
	// 			"type":                "dateFile",
	// 			"filename":            "./temp/test",
	// 			"maxLogSize":          130,
	// 			"keepFileExt":         false,
	// 			"channelSize":         160,
	// 			"internalChannelSize": 120,
	// 			"mode":                493,
	// 			"flags":               511,
	// 			"compress":            true,
	// 			"compressMode":        "default",
	// 			"fileNameSep":         ".test.",
	// 			"layout": map[string]any{
	// 				"type":    "pattern",
	// 				"pattern": "%[%d %c%] ~ %m",
	// 			},
	// 		},
	// 	},
	// 	"categories": map[string]any{
	// 		"cheese": map[string]any{
	// 			"appenders": []string{
	// 				"out",
	// 			},
	// 			"level": "error",
	// 		},
	// 		"default": map[string]any{
	// 			"appenders": []string{
	// 				"out",
	// 			},
	// 			"level": "all",
	// 		},
	// 	},
	// })
	Configure("./example/logger-config.json")
	l := GetLogger("app")
	l.Info("info")
	l.Fatal("a fatal error")

	logger := GetLogger()
	logger.Info("info")
	logger.Fatal("a fatal error")
	logger.Info("5")
	logger.Info("4")
	logger.Info("3")
	logger.Info("2")
	logger.Info("1")
	logger.Info("hot file")
	for i := 0; i < 10000; i++ {
		logger.Trace("TRACE log")
		logger.Debug("DEBUG log")
		logger.Info("INFO Log")
		logger.Warn("WARN Log")
		logger.Error("ERROR Log")
		logger.Fatal("FATAL log")
	}
	logger.Info("Before Shutdown")
	Shutdown()
	logger.Info("After Shutdown")
}
