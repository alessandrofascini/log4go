package appenders

import (
	"container/heap"
	"errors"
	"io"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/alessandrofascini/log4go/pkg"

	dsa "github.com/alessandrofascini/log4go/internal/pkg"

	errorhelper "github.com/alessandrofascini/log4go/internal/errors"
)

func newDateFileManager(config *FileConfig) IFileManagerRotation {
	if config.maxLogSize == 0 {
		// first case
		cmp := NewDateComparator(config.filename, config.fileNameSep, config.keepFileExt, config.compress, config.pattern)
		return &DateFileManagerI{
			filename:     config.filename,
			flag:         config.flags,
			mode:         config.mode,
			maxLogSize:   config.maxLogSize,
			backups:      uint(config.backups),
			compress:     config.compress,
			pattern:      config.pattern,
			comparator:   &cmp,
			compressMode: config.compressMode,
			channelSize:  config.channelSize,
		}
	}
	// second case
	cmp := NewDateComparatorII(config.filename, config.fileNameSep, config.keepFileExt, config.compress, config.pattern)
	return &DateFileManagerII{
		filename:     config.filename,
		flag:         config.flags,
		mode:         config.mode,
		maxLogSize:   config.maxLogSize,
		backups:      uint(config.backups),
		compress:     config.compress,
		pattern:      config.pattern,
		comparator:   &cmp,
		compressMode: config.compressMode,
		channelSize:  config.channelSize,
	}
}

// second case: Date Rolling File Appender

type DateFileManagerI struct {
	filename     string
	flag         int
	mode         os.FileMode
	maxLogSize   int64
	backups      uint
	compress     bool
	pattern      string
	comparator   *DateComparator
	compressMode int
	channelSize  int
}

func (dfm *DateFileManagerI) getHotFile() string {
	return dfm.filename
}

func (dfm *DateFileManagerI) getFlags() int {
	return dfm.flag
}

func (dfm *DateFileManagerI) getMode() os.FileMode {
	return dfm.mode
}

func (dfm *DateFileManagerI) isRequiredRotation() bool {
	file, err := os.Open(dfm.filename)
	if err != nil {
		var e *os.PathError
		if errors.As(err, &e) && e.Err != syscall.ENOENT {
			errorhelper.WriteErrorf("cannot read %q, %+v\n", dfm.filename, err)
		}
		return false
	}
	defer func() {
		if err = file.Close(); err != nil {
			errorhelper.WriteError(err)
		}
	}()
	stats, err := file.Stat()
	if err != nil {
		errorhelper.WriteErrorf("cannot read stats of %q, %+v\n", dfm.filename, err)
		return false
	}
	// #1 date file
	return stats.ModTime().Format(dfm.pattern) != pkg.GetDateByFormat(dfm.pattern)
}

// rollNow here implementation date changed rotation
func (dfm *DateFileManagerI) rollNow() {
	path, hotFilename, _ := getPathHotFileExt(&dfm.filename)

	filesChannel, _ := readDateFiles(path, dfm.channelSize)
	chDeleteFiles := dateMatchOrDelete(filesChannel, dfm.comparator, int(dfm.backups), dfm.channelSize)

	done := deleteFiles(chDeleteFiles, path)

	if dfm.compress {
		if err := compress(path+hotFilename, path+dfm.comparator.Replace(0), dfm.compressMode); err != nil {
			errorhelper.WriteError(err)
		}
		// Now I can Rotate the hot file
	} else if err := os.Rename(path+hotFilename, path+dfm.comparator.Replace(0)); err != nil {
		errorhelper.WriteError(err)
	}
	<-done
}

func readDateFiles(path string, channelSize int) (chan string, error) {
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

func dateMatchOrDelete(files chan string, comp *DateComparator, numBackups, channelSize int) chan string {
	chStr := make(chan string, channelSize)
	go func() {
		for file := range files {
			m := comp.Match(file)
			if m >= numBackups {
				chStr <- file
			}
		}
		close(chStr)
	}()
	return chStr
}

// second case: Date Rolling File Appender & File Appender

type DateFileManagerII struct {
	filename     string
	flag         int
	mode         os.FileMode
	maxLogSize   int64
	backups      uint
	compress     bool
	pattern      string
	comparator   *DateComparatorII
	compressMode int
	channelSize  int
}

func (dfm *DateFileManagerII) getHotFile() string {
	return dfm.filename
}

func (dfm *DateFileManagerII) getFlags() int {
	return dfm.flag
}

func (dfm *DateFileManagerII) getMode() os.FileMode {
	return dfm.mode
}

func (dfm *DateFileManagerII) isRequiredRotation() bool {
	/* Two type of Rotation:
	1) file dimensions
	2) date changed
	*/
	file, err := os.Open(dfm.filename)
	if err != nil {
		var e *os.PathError
		if errors.As(err, &e) && e.Err != syscall.ENOENT {
			errorhelper.WriteErrorf("cannot read %q, %+v\n", dfm.filename, err)
		}
		return false
	}
	defer func() {
		if err = file.Close(); err != nil {
			errorhelper.WriteError(err)
		}
	}()
	stats, err := file.Stat()
	if err != nil {
		errorhelper.WriteErrorf("cannot read stats of %q, %+v\n", dfm.filename, err)
		return false
	}
	// #1 file dimensions || #2 date file
	return (dfm.maxLogSize > 0 && dfm.maxLogSize < stats.Size()) || stats.ModTime().Format(dfm.pattern) != pkg.GetDateByFormat(dfm.pattern)
}

// rollNow here implementation date changed rotation
func (dfm *DateFileManagerII) rollNow() {
	path, hotFilename, _ := getPathHotFileExt(&dfm.filename)

	filesChannel, _ := readDateFiles(path, dfm.channelSize)
	chMatchFiles, chDeleteFiles := dfm.matchOrDelete(filesChannel, dfm.comparator, int64(dfm.backups))

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		deleteFilesSync(chDeleteFiles, path)
		wg.Done()
	}()

	// ------
	heapArr := make([]dsa.IntHeap, dfm.backups)
	for matchFile := range chMatchFiles {
		heap.Push(&heapArr[matchFile[0]], int(matchFile[1]))
	}

	fileCounter := 0
	i := 0
	for ; i < len(heapArr) && fileCounter < int(dfm.backups); i++ {
		fileCounter += heapArr[i].Len()
	}
	i--
	// -------

	if i == 0 {
		now := time.Now()
		midnightToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		dfm.fixBackups(&heapArr[0], midnightToday, fileCounter-int(dfm.backups)+1, &path)
	} else {
		wg.Add(1)
		go func() {
			now := time.Now()
			midnightToday := time.Date(now.Year(), now.Month(), now.Day()-i, 0, 0, 0, 0, time.UTC)
			dfm.fixBackups(&heapArr[i], midnightToday, fileCounter-int(dfm.backups), &path)
			wg.Done()
		}()
	}

	wg.Add(2)
	go func() {
		dfm.rotateToday(&heapArr[0], path, hotFilename)
		wg.Done()
	}()

	go func() {
		dfm.deleteExtraFiles(heapArr, i+1, &path)
		wg.Done()
	}()

	wg.Wait()
}

func (dfm *DateFileManagerII) matchOrDelete(files chan string, comp *DateComparatorII, numBackups int64) (chan [2]int64, chan string) {
	chInt := make(chan [2]int64, dfm.channelSize)
	chStr := make(chan string, dfm.channelSize)
	go func() {
		for file := range files {
			m, off := comp.Match(file)
			if m > -1 && m < numBackups && off != 0 { // [0, numBackups)
				chInt <- [2]int64{m, int64(off)}
			} else if m != -1 {
				chStr <- file
			}
		}
		close(chInt)
		close(chStr)
	}()
	return chInt, chStr
}

func (dfm *DateFileManagerII) rotateToday(h *dsa.IntHeap, path, hotFilename string) {
	// current day rotation
	// presents files
	presentFiles := make([]bool, heap.Pop(h).(int))

	if len(presentFiles) > 0 {
		presentFiles[len(presentFiles)-1] = true
	}

	for h.Len() > 0 {
		presentFiles[heap.Pop(h).(int)-1] = true
	}

	// get first renameable file
	firstRenameableFile := 0
	for firstRenameableFile < len(presentFiles) && presentFiles[firstRenameableFile] {
		firstRenameableFile++
	}
	next := path + dfm.comparator.Replace(firstRenameableFile+1)
	for ; firstRenameableFile > 0; firstRenameableFile-- {
		curr := path + dfm.comparator.Replace(firstRenameableFile)
		if err := os.Rename(curr, next); err != nil {
			errorhelper.WriteError(err)
		}
		next = curr
	}

	if dfm.compress {
		err := compress(path+hotFilename, next, dfm.compressMode)
		if err != nil {
			errorhelper.WriteError(err)
		}
		return
	}

	// Now I can Rotate the hot file
	if err := os.Rename(path+hotFilename, next); err != nil {
		errorhelper.WriteError(err)
	}
}

func (dfm *DateFileManagerII) deleteExtraFiles(heapArr []dsa.IntHeap, j int, path *string) {
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day()-j, 0, 0, 0, 0, time.UTC)
	for ; j < len(heapArr); j++ {
		for _, subLog := range heapArr[j] {
			if err := os.Remove(*path + dfm.comparator.ReplaceWithDifferentDate(subLog, date)); err != nil {
				errorhelper.WriteError(err)
			}
		}
		date = date.Add(-1 * 24 * time.Hour) // add 1 day
	}
}

func (dfm *DateFileManagerII) fixBackups(h *dsa.IntHeap, date time.Time, countUntil int, path *string) {
	for h.Len() > 0 && countUntil != 0 {
		v := heap.Pop(h).(int)
		if err := os.Remove(*path + dfm.comparator.ReplaceWithDifferentDate(v, date)); err != nil {
			errorhelper.WriteError(err)
		}
		countUntil--
	}
}
