package tests

import (
	"container/heap"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alessandrofascini/log4go/pkg"
)

// DateComparator struct
type DateComparator struct {
	sep         string
	pattern     string
	datePattern string
}

func NewDateComparator(fullPath string, sep string, keepFileExt, isCompress bool, datePattern string) DateComparator {
	_, fullFileName := filepath.Split(fullPath)
	ext := filepath.Ext(fullFileName)
	filename := fullFileName[0 : len(fullFileName)-len(ext)]
	p := strings.Builder{}
	if keepFileExt {
		p.WriteString(filename)
		p.WriteString(string(SpecialChar))
		p.WriteString(ext)
	} else {
		p.WriteString(fullFileName)
		p.WriteString(string(SpecialChar))
	}
	if isCompress {
		p.WriteString(CompressExt)
	}
	return DateComparator{
		sep:         sep,
		pattern:     p.String(),
		datePattern: datePattern,
	}
}

const days = 24 * 60 * 60 * 1000

// Match return the difference between now and the date saved
// -1
func (d *DateComparator) Match(s string) int {
	i := 0
	for ; i < len(d.pattern) && d.pattern[i] != SpecialChar; i++ {
		if d.pattern[i] != s[i] {
			return -1
		}
	}
	j := len(s) - 1
	w := len(d.pattern) - 1
	for ; w > -1 && d.pattern[w] != SpecialChar; w-- {
		if d.pattern[w] != s[j] {
			return -1
		}
		j--
	}
	j++         // j is not included
	if i == j { // We found our hot file
		return -1
	}
	// Now indexing
	if j-i < len(d.sep) {
		return -1
	}
	w = 0
	for i < j && w < len(d.sep) {
		if s[i] != d.sep[w] {
			return -1
		}
		i++
		w++
	}
	if i == j { // This is useful for date file appender
		return -1
	}
	// between 'i' and 'j' we will find our date
	if i+len(d.datePattern) > j {
		return -1
	}
	v := s[i : i+len(d.datePattern)]
	parsed, err := time.Parse(d.datePattern, v)
	if err != nil {
		return -1
	}
	now := time.Now()
	// in days
	diff := (now.UnixMilli() - parsed.UnixMilli()) / days
	if diff < 0 {
		return -1
	}
	return int(diff)
}

func (d *DateComparator) Replace(i int) string {
	date := fmt.Sprintf("%s%s", d.sep, pkg.GetDateByFormat(d.datePattern))
	if i != 0 {
		date = fmt.Sprintf("%s%s%d", date, d.sep, i)
	}
	return strings.Replace(d.pattern, string(SpecialChar), date, -1)
}

// case I, a file is using only dates

func DateFileRotationIICaseI(filename string, backups uint) {
	path, hotFilename, _ := getPathHotFileExt(filename)
	comp := NewDateComparator(hotFilename, ".", false, false, "2006-01-02")

	filesChannel, _ := readDateFiles(path)
	chDeleteFiles := dateMatchOrDelete(filesChannel, &comp, int(backups))

	done := deleteFiles(chDeleteFiles, path)

	// rename hot file
	if err := os.Rename(path+hotFilename, path+comp.Replace(0)); err != nil {
		WriteError(err)
	}

	<-done
}

func readDateFiles(path string) (chan string, error) {
	dir, err := os.Open(path)
	channel := make(chan string, ChannelSize)
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

func dateMatchOrDelete(files chan string, comp *DateComparator, numBackups int) chan string {
	chStr := make(chan string, ChannelSize)
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

// Case II

type DateComparatorII struct {
	sep         string
	pattern     string
	datePattern string
}

func NewDateComparatorII(fullPath string, sep string, keepFileExt, isCompress bool, datePattern string) DateComparatorII {
	_, fullFileName := filepath.Split(fullPath)
	ext := filepath.Ext(fullFileName)
	filename := fullFileName[0 : len(fullFileName)-len(ext)]
	p := strings.Builder{}
	if keepFileExt {
		p.WriteString(filename)
		p.WriteString(string(SpecialChar))
		p.WriteString(ext)
	} else {
		p.WriteString(fullFileName)
		p.WriteString(string(SpecialChar))
	}
	if isCompress {
		p.WriteString(CompressExt)
	}
	return DateComparatorII{
		sep:         sep,
		pattern:     p.String(),
		datePattern: datePattern,
	}
}

// Match return the difference between now and the date saved
// -1
func (d *DateComparatorII) Match(s string) (int64, int) {
	i := 0
	for ; i < len(d.pattern) && d.pattern[i] != SpecialChar; i++ {
		if d.pattern[i] != s[i] {
			return -1, -1
		}
	}
	j := len(s) - 1
	w := len(d.pattern) - 1
	for ; w > -1 && d.pattern[w] != SpecialChar; w-- {
		if d.pattern[w] != s[j] {
			return -1, -1
		}
		j--
	}
	j++         // j is not included
	if i == j { // We found our hot file
		return -1, -1
	}
	// Now indexing
	if j-i < len(d.sep) {
		return -1, -1
	}
	w = 0
	for i < j && w < len(d.sep) {
		if s[i] != d.sep[w] {
			return -1, -1
		}
		i++
		w++
	}
	if i == j { // This is useful for date file appender
		return -1, -1
	}
	// between 'i' and 'j' we will find our date
	if i+len(d.datePattern) > j {
		return -1, -1
	}
	dateAndNumber := strings.Split(s[i:j], d.sep)
	if len(dateAndNumber) == 0 {
		return -1, -1
	}

	parsed, err := time.Parse(d.datePattern, dateAndNumber[0])
	if err != nil {
		return -1, -1
	}

	now := time.Now()
	midnightToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// in days
	diff := (midnightToday.UnixMilli() - parsed.UnixMilli()) / days
	if diff < 0 {
		return -1, -1
	}

	if len(dateAndNumber) < 2 {
		return diff, 0
	}

	// calculate offset => ${specialChar}${number}
	offset, err := strconv.Atoi(dateAndNumber[1])
	if err != nil {
		return -1, -1
	}
	return diff, offset
}

func (d *DateComparatorII) Replace(i int) string {
	date := fmt.Sprintf("%s%s", d.sep, pkg.GetDateByFormat(d.datePattern))
	if i != 0 {
		date = fmt.Sprintf("%s%s%d", date, d.sep, i)
	}
	return strings.Replace(d.pattern, string(SpecialChar), date, -1)
}

func (d *DateComparatorII) ReplaceWithDifferentDate(i int, date time.Time) string {
	dateStr := fmt.Sprintf("%s%s", d.sep, date.Format(d.datePattern))
	if i != 0 {
		dateStr = fmt.Sprintf("%s%s%d", dateStr, d.sep, i)
	}
	return strings.Replace(d.pattern, string(SpecialChar), dateStr, -1)
}

func DateFileRotationIICaseII(filename string, backups uint) {
	path, hotFilename, _ := getPathHotFileExt(filename)
	comp := NewDateComparatorII(hotFilename, ".", false, false, "2006-01-02")

	filesChannel, _ := readDateFiles(path)
	chMatchFiles, chDeleteFiles := dateMatchOrDeleteII(filesChannel, &comp, int64(backups))

	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		deleteFilesSync(chDeleteFiles, path)
	}()

	// ------
	heapArr := make([]IntHeap, backups)
	for matchFile := range chMatchFiles {
		heap.Push(&heapArr[matchFile[0]], int(matchFile[1]))
	}

	fileCounter := 0
	i := 0
	for ; i < len(heapArr) && fileCounter <= int(backups); i++ {
		fileCounter += heapArr[i].Len()
	}
	i--
	// -------
	fmt.Println(i)
	if i == 0 {
		now := time.Now()
		midnightToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		fixBackups(&heapArr[0], midnightToday, fileCounter-int(backups), &comp, path)
	} else {
		go func() {
			wg.Add(1)
			defer wg.Done()
			now := time.Now()
			midnightToday := time.Date(now.Year(), now.Month(), now.Day()-i, 0, 0, 0, 0, time.UTC)
			fixBackups(&heapArr[i], midnightToday, fileCounter-int(backups), &comp, path)
		}()
	}

	go func() {
		wg.Add(1)
		defer wg.Done()
		rotateToday(&heapArr[0], &comp, path, hotFilename)
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		deleteExtraFiles(heapArr, i+1, path, &comp)
	}()

	wg.Wait()
}

func dateMatchOrDeleteII(files chan string, comp *DateComparatorII, numBackups int64) (chan [2]int64, chan string) {
	chInt := make(chan [2]int64, ChannelSize)
	chStr := make(chan string, ChannelSize)
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

func rotateToday(h *IntHeap, comp *DateComparatorII, path, hotFilename string) {
	// current day rotation
	// presents files
	presentFiles := make([]bool, heap.Pop(h).(int))
	for h.Len() > 0 {
		presentFiles[heap.Pop(h).(int)] = true
	}

	// get first renameable file
	firstRenameableFile := 1
	for firstRenameableFile < len(presentFiles) && presentFiles[firstRenameableFile] {
		firstRenameableFile++
	}
	next := path + comp.Replace(firstRenameableFile+1)
	for ; firstRenameableFile > 0; firstRenameableFile-- {
		curr := path + comp.Replace(firstRenameableFile)
		if err := os.Rename(curr, next); err != nil {
			WriteError(err)
		}
		next = curr
	}
	// Now I can Rotate the hot file
	if err := os.Rename(path+hotFilename, next); err != nil {
		WriteError(err)
	}
}

func deleteExtraFiles(heapArr []IntHeap, j int, path string, comp *DateComparatorII) {
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day()-j, 0, 0, 0, 0, time.UTC)
	for ; j < len(heapArr); j++ {
		for _, subLog := range heapArr[j] {
			if err := os.Remove(path + comp.ReplaceWithDifferentDate(subLog, date)); err != nil {
				WriteError(err)
			}
		}
		date = date.Add(-1 * 24 * time.Hour) // add 1 day
	}
}

func fixBackups(h *IntHeap, date time.Time, countUntil int, comp *DateComparatorII, path string) {
	for h.Len() > 0 && countUntil != 0 {
		v := heap.Pop(h).(int)
		if err := os.Remove(path + comp.ReplaceWithDifferentDate(v, date)); err != nil {
			WriteError(err)
		}
		countUntil--
	}
}

func deleteFilesSync(chFiles chan string, path string) {
	for file := range chFiles {
		if err := os.Remove(filepath.Join(path, file)); err != nil {
			WriteError(err)
		}
	}
}
