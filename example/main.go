package main

import (
	"fmt"
	"log"
	"time"

	"github.com/alessandrofascini/log4go/pkg"

	log4go "github.com/alessandrofascini/log4go"
)

func main() {
	start := time.Now()
	customAppender := func(config *pkg.AppenderConfig) pkg.Appender {
		return func(le pkg.LoggingEvent) {
			// do nothing
		}
	}
	log4go.AddAppender("custom", &customAppender)
	if err := log4go.Configure("logger-config.json"); err != nil {
		log.Fatal(err)
	}
	//log4go.Configure(map[string]any{
	//	"appenders": map[string]any{
	//		"out": map[string]any{
	//			"type":                "dateFile",
	//			"filename":            "./temp/test",
	//			"maxLogSize":          130,
	//			"keepFileExt":         false,
	//			"channelSize":         160,
	//			"internalChannelSize": 120,
	//			"mode":                0755,
	//			"flags":               os.O_APPEND | os.O_CREATE | os.O_WRONLY,
	//			"compress":            true,
	//			"compressMode":        "default",
	//			"fileNameSep":         ".test.",
	//			"layout": map[string]any{
	//				"type":    "pattern",
	//				"pattern": "%[%d %c%] ~ %m",
	//			},
	//		},
	//		"app": map[string]any{
	//			"type":         "stdout",
	//			"filename":     "./temp/app.txt",
	//			"maxLogSize":   2500,
	//			"keepFileExt":  true,
	//			"compress":     false,
	//			"compressMode": "default",
	//			"backups":      5,
	//			"layout": map[string]any{
	//				"type":    "pattern",
	//				"pattern": "%[%d %c%] ~ %m",
	//			},
	//			"channelSize": 20,
	//		},
	//	},
	//	"categories": map[string]any{
	//		"default": map[string]any{
	//			"appenders": []string{
	//				"out",
	//				"app",
	//			},
	//			"level": "all",
	//		},
	//		"app": map[string]any{
	//			"appenders": []string{
	//				"app",
	//			},
	//			"level": "all",
	//		},
	//	},
	//	"configuration": map[string]any{
	//		"appenderChannelSize":         160,
	//		"internalAppenderChannelSize": 110,
	//		"createFolderPerm":            0750,
	//	},
	//})
	logger := log4go.GetLogger()
	appLogger := log4go.GetLogger()
	appLogger.Info("ciao")
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
	appLogger.Info("prima della fine")
	log4go.Shutdown()
	logger.Info("After Shutdown")
	fmt.Println("finished after ", time.Now().UnixMilli()-start.UnixMilli())
	fmt.Println(time.Now().Format("2006-01-02T15:04:05.000"))
}
