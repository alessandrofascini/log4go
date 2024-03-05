package appenders

import (
	"os"
	"path/filepath"

	"github.com/alessandrofascini/log4go/pkg"

	errorshelper "github.com/alessandrofascini/log4go/internal/errors"

	configuration "github.com/alessandrofascini/log4go/internal/config"
)

/*
IFileManagerRotation
This interface defines methods for file's manager appenders.
*/
type IFileManagerRotation interface {
	getHotFile() string
	getFlags() int
	getMode() os.FileMode
	isRequiredRotation() bool
	rollNow()
}

/*
fileAppenderFactory
This internal function allow us to reuse the same code for file's appenders
*/
func fileAppenderFactory(config pkg.AppenderConfig, layout *pkg.Layout, managerFactory FileManagerFactory) pkg.Appender {
	manager := managerFactory(newFileConfig(config))
	return func(event pkg.LoggingEvent) {
		// #1 	Check if I need to rotate
		if manager.isRequiredRotation() {
			// #1.1 If true rotate
			manager.rollNow()
		}
		// #2   Write content on hot_file
		dir := filepath.Dir(manager.getHotFile())
		_, err := os.Stat(dir)
		if err != nil {
			err = os.Mkdir(dir, configuration.GetCreateFolderPerm())
			if err != nil && !os.IsExist(err) {
				panic(err)
			}
		}
		file, err := os.OpenFile(manager.getHotFile(), manager.getFlags(), manager.getMode())
		if err != nil {
			errorshelper.WriteError(err)
			return
		}
		if _, err = file.WriteString((*layout)(event)); err != nil {
			errorshelper.WriteError(err)
			return
		}
		if err = file.Close(); err != nil {
			errorshelper.WriteError(err)
			return
		}
	}
}

/*
FileSyncAppender is a function that creates an Appender of type

	FileSync appender config
	{
		"type": "fileSync" (string),
		"filename": string
		"maxLogSize": integer | string (OPTIONAL: default 0).
			The maximum size (in bytes) for the log file.
			If not specified or 0, then no log rolling will happen.
			maxLogSize can also accept string with the size suffixes: K, M, G such as 1K, 1M, 1G.
		"backups": integer (OPTIONAL: default 5).
			The number of old log files to keep during log rolling (excluding the hot file).
		"layout": a Layout. (already configured)
		"compressMode": how to compress the file
	}
*/
func FileSyncAppender(config pkg.AppenderConfig, layout *pkg.Layout) pkg.Appender {
	return fileAppenderFactory(config, layout, newFileSyncManager)
}

/*
FileAppender is a function that creates an Appender of type

	File appender config
	{
		"type": "fileSync" (string),
		"filename": string
		"maxLogSize": integer | string (OPTIONAL: default 0).
			The maximum size (in bytes) for the log file.
			If not specified or 0, then no log rolling will happen.
			maxLogSize can also accept string with the size suffixes: K, M, G such as 1K, 1M, 1G.
		"backups": integer (OPTIONAL: default 5).
			The number of old log files to keep during log rolling (excluding the hot file).
		"layout": a Layout. (already configured)
		"compress": boolean (default false)
			 compress the backup files using gzip (backup files will have .gz extension).
		"compressMode": int | string
			how to compress the file
		"keepFileExt": boolean (default false)
			preserve the file extension when rotating log files (file.log becomes file.1.log instead of file.log.1).
		"fileNameSep": string (default ".")
	}
*/
func FileAppender(config pkg.AppenderConfig, layout *pkg.Layout) pkg.Appender {
	return fileAppenderFactory(config, layout, newFileManager)
}

/*
DateFileAppender is a function that create an DateFile Appender

	File appender config
	{
		"type": "fileSync" (string),
		"filename": string
		"pattern": string (optional, defaults to 2006-01-02).
			The pattern to use to determine when to roll the logs.
			The pattern use the standard golang time format style.
		"maxLogSize": integer | string (OPTIONAL: default 0).
			The maximum size (in bytes) for the log file.
			If not specified or 0, then no log rolling will happen.
			maxLogSize can also accept string with the size suffixes: K, M, G such as 1K, 1M, 1G.
		"backups": integer (OPTIONAL: default 5).
			The number of old log files to keep during log rolling (excluding the hot file).
		"layout": a Layout. (already configured)
		"compress": boolean (default false)
			 compress the backup files using gzip (backup files will have .gz extension).
		"compressMode": int | string
			how to compress the file.
		"keepFileExt": boolean (default false)
			preserve the file extension when rotating log files (file.log becomes file.1.log instead of file.log.1).
		"fileNameSep": string (default ".")
	}
*/
func DateFileAppender(config pkg.AppenderConfig, layout *pkg.Layout) pkg.Appender {
	return fileAppenderFactory(config, layout, newDateFileManager)
}
