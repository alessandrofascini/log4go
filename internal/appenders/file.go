package appenders

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	errorshelper "github.com/alessandrofascini/log4go/internal/errors"

	configuration "github.com/alessandrofascini/log4go/internal/config"
)

// Functions
func newFileManager(config *FileConfig) IFileManagerRotation {
	c := NewComparator(config.filename, config.fileNameSep, config.keepFileExt, config.compress)
	return &FileManager{
		filename:    config.filename,
		flag:        config.flags,
		mode:        config.mode,
		maxLogSize:  config.maxLogSize,
		backups:     uint(config.backups),
		compress:    config.compress,
		comparator:  c,
		channelSize: configuration.GetDefaultAppenderChannelSize(),
	}
}

type FileManager struct {
	filename     string
	flag         int
	mode         os.FileMode
	maxLogSize   int64
	backups      uint
	compress     bool
	compressMode int
	comparator   *Comparator
	channelSize  int
}

func (fm *FileManager) getFlags() int {
	return fm.flag
}

func (fm *FileManager) getMode() os.FileMode {
	return fm.mode
}

func (fm *FileManager) getHotFile() string {
	return fm.filename
}

func (fm *FileManager) isRequiredRotation() bool {
	if fm.maxLogSize == 0 {
		return false
	}
	file, err := os.Open(fm.filename)
	var e *os.PathError
	if errors.As(err, &e) && e.Err != syscall.ENOENT {
		errorshelper.WriteErrorf("cannot read %q, %+v\n", fm.filename, err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			errorshelper.WriteError(err)
		}
	}()
	stats, err := file.Stat()
	if err != nil {
		errorshelper.WriteErrorf("cannot read stats of %q, %+v", fm.filename, err)
		return false
	}
	return fm.maxLogSize < stats.Size()
}

func (fm *FileManager) rollNow() {
	path, hotFilename, _ := getPathHotFileExt(&fm.filename)

	// read files
	filesChannel, _ := readFiles(path, fm.channelSize)
	// match file or delete it
	chValidFiles, chDeleteFiles := matchOrDeleteV(filesChannel, fm.comparator, int(fm.backups), fm.channelSize)

	// delete files in a goroutine (WARNING)
	done := deleteFiles(chDeleteFiles, path)

	// presents files
	presentFiles := make([]bool, int(fm.backups))
	for f := range chValidFiles {
		presentFiles[f] = true
	}

	// get first renameable file
	firstRenameableFile := 1
	for firstRenameableFile < len(presentFiles) && presentFiles[firstRenameableFile] {
		firstRenameableFile++
	}
	next := path + fm.comparator.Replace(firstRenameableFile)
	for firstRenameableFile = firstRenameableFile - 1; firstRenameableFile > 0; firstRenameableFile-- {
		curr := path + fm.comparator.Replace(firstRenameableFile)
		if err := os.Rename(curr, next); err != nil {
			errorshelper.WriteError(err)
		}
		next = curr
	}
	if fm.compress {
		err := compress(path+hotFilename, next, fm.compressMode)
		if err != nil {
			errorshelper.WriteError(err)
		}
		// Now I can Rotate the hot file
	} else if err := os.Rename(path+hotFilename, next); err != nil {
		errorshelper.WriteError(err)
	}
	<-done
}

func getPathHotFileExt(filename *string) (path, hotFilename, ext string) {
	path, hotFilename = filepath.Split(*filename)
	ext = filepath.Ext(hotFilename)[1:]
	return
}

func readFiles(path string, channelSize int) (chan string, error) {
	dir, err := os.Open(path)
	channel := make(chan string, channelSize)
	if err != nil {
		close(channel)
		return channel, err
	}
	go func() {
		incrementer := 32
		var e error
		var n []string
		var wg sync.WaitGroup
		for e == nil || e != io.EOF {
			n, e = dir.Readdirnames(incrementer)
			incrementer = incrementer<<1 + 1
			wg.Add(1)
			go func(n []string) {
				for _, v := range n {
					channel <- v
				}
				wg.Done()
			}(n)
		}
		wg.Wait()
		close(channel)
	}()
	return channel, nil
}

func matchOrDeleteV(files chan string, comp *Comparator, numBackups, channelSize int) (chan int, chan string) {
	chInt := make(chan int, channelSize)
	chStr := make(chan string, channelSize)
	go func() {
		for file := range files {
			m := comp.Match(file)
			if m > -1 {
				if m < numBackups {
					chInt <- m
				} else {
					chStr <- file
				}
			}
		}
		close(chInt)
		close(chStr)
	}()
	return chInt, chStr
}

func deleteFilesSync(chFiles chan string, path string) {
	for file := range chFiles {
		if err := os.Remove(filepath.Join(path, file)); err != nil {
			errorshelper.WriteError(err)
		}
	}
}

func deleteFiles(chFiles chan string, path string) chan interface{} {
	done := make(chan interface{})
	go func(chFiles chan string, path string) {
		deleteFilesSync(chFiles, path)
		close(done)
	}(chFiles, path)
	return done
}
