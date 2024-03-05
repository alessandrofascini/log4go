package appenders

import (
	"os"
)

// Constants
// Variables
// Functions

func newFileSyncManager(config *FileConfig) IFileManagerRotation {
	const (
		keepFileExt = true
		compress    = true
		sep         = "."
	)
	c := NewComparator(config.filename, sep, keepFileExt, compress)
	return &FileManager{
		filename:   config.filename,
		flag:       config.flags | os.O_SYNC,
		mode:       config.mode,
		maxLogSize: config.maxLogSize,
		backups:    uint(config.backups),
		compress:   compress,
		comparator: c,
	}
}
