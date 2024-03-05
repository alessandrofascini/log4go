package tests

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

//const hotFilename = "mon.log"

/*
Design Pattern Bridge
For more details, go to https://refactoring.guru/design-patterns/bridge
*/
//type fileManager struct {
//	filename string
//	rotator  fileRotation
//}
//
//type fileRotation struct {
//	maxLogSize uint64
//	backups    uint
//	compress   bool
//	matcher    string
//	archiver   archiveHotFile
//}
//
//type archiveHotFile func()
//
//func appenderFactory() *fileManager {
//	return &fileManager{
//		filename: filepath.Join(path, hotFilename),
//		rotator:  fileRotation{},
//	}
//}

func filenameFactory(filename string) (path, hotFilename, ext string) {
	path, hotFilename = filepath.Split(filename)
	ext = filepath.Ext(hotFilename)
	return
}

func FileRotation(filename string, backups uint) {
	path, hotFilename, _ := filenameFactory(filename)
	// Open Directory
	dir, _ := os.Open(path)
	// Get file names
	names, _ := dir.Readdirnames(0)
	// Matcher (mon.log.1)
	pattern := fmt.Sprintf("%s.*", hotFilename)

	// Filter
	var matches []string
	for _, name := range names {
		if v, _ := filepath.Match(pattern, name); v {
			matches = append(matches, name)
		}
	}

	// Sort
	sort.Strings(matches)

	// Truncation
	matches = matches[0:min(backups, uint(len(matches)))]

	// Reverse
	reverseArray(matches)

	//fmt.Println(matches)
	// Update Filenames
	//for _, name := range matches {
	//	ext := filepath.Ext(name)
	//	i, _ := strconv.Atoi(ext[1:])
	//	n := name[0 : len(name)-len(ext)]
	//	_ = os.Rename(filepath.Join(path, name), filepath.Join(path, fmt.Sprintf("%s.%d", n, i+1)))
	//}

	// update hot file
	//_ = os.Rename(filepath.Join(path, hotFilename), filepath.Join(path, fmt.Sprintf("%s.%d", hotFilename, 1)))
}

func reverseArray(arr []string) {
	left := 0
	right := len(arr) - 1
	for left < right {
		t := arr[left]
		arr[left] = arr[right]
		arr[right] = t
		left++
		right--
	}
}

func min(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}

const SpecialChar = ';'
const CompressExt = ".gz"

type Comparator struct {
	pattern string
	sep     string
}

func NewComparator(fullPath string, sep string, keepFileExt, isCompress bool) *Comparator {
	_, fullFileName := filepath.Split(fullPath)
	ext := filepath.Ext(fullFileName)
	filename := fullFileName[0 : len(fullFileName)-len(ext)]
	var p []string
	if keepFileExt {
		p = append(p, filename, string(SpecialChar), ext)
	} else {
		p = append(p, fullFileName, string(SpecialChar))
	}
	if isCompress {
		p = append(p, CompressExt)
	}
	return &Comparator{
		pattern: strings.Join(p, ""),
		sep:     sep,
	}
}

/*
Match
This method return -1 if no matching;
0 if it is the hot file
greater than 0 if there is a matching
*/
func (c *Comparator) Match(s string) int {
	i := 0
	for ; i < len(c.pattern) && c.pattern[i] != SpecialChar; i++ {
		if c.pattern[i] != s[i] {
			return -1
		}
	}
	j := len(s) - 1
	w := len(c.pattern) - 1
	for ; w > -1 && c.pattern[w] != SpecialChar; w-- {
		if c.pattern[w] != s[j] {
			return -1
		}
		j--
	}
	j++         // j is not included
	if i == j { // We found our hot file
		return 0
	}
	// Now indexing
	if j-i < len(c.sep) {
		return -1
	}
	w = 0
	for i < j && w < len(c.sep) {
		if s[i] != c.sep[w] {
			return -1
		}
		i++
		w++
	}
	if i == j { // This is useful for date file appender
		return 0
	}
	// between 'i' and 'j' we will find our number
	res := 0
	for i < j {
		res = res<<1 + res<<3
		b := int(s[i] - '0')
		res += b
		i++
	}
	return res
}

func (c *Comparator) Replace(i int) string {
	// Is better to create a custom function to do this?
	newStr := fmt.Sprintf("%s%d", c.sep, i)
	return strings.Replace(c.pattern, string(SpecialChar), newStr, -1)
}

func getPathHotFileExt(filename string) (path, hotFilename, ext string) {
	path, hotFilename = filepath.Split(filename)
	ext = filepath.Ext(hotFilename)[1:]
	return
}

func FileRotationIII(filename string, backups uint) {
	// Opening folder
	path, hotFilename, _ := getPathHotFileExt(filename)

	// Open Directory
	dir, _ := os.Open(path)
	defer dir.Close()

	comparator := NewComparator(hotFilename, ".", false, false)

	B := int(backups + 1)
	fileBackups := make([]string, B)
	incrementer := 32
	var n []string
	var e error
	var wg sync.WaitGroup
	for e == nil || e != io.EOF {
		n, e = dir.Readdirnames(incrementer)
		incrementer = incrementer<<1 + 1
		// if matches
		wg.Add(1)
		go func(n []string) {
			defer wg.Done()
			for _, v := range n {
				// ## MATCH ##
				m := comparator.Match(v)
				if m == -1 { // IS MATCHING and isn't overflow
					continue
				}
				if m < B {
					fileBackups[m] = v
				} else if err := os.Remove(filepath.Join(path, v)); err != nil {
					fmt.Println("deleted", filepath.Join(path, v))
					// overflow => delete it
					WriteError(err)
				}
			}
		}(n)
	}
	wg.Wait()
	//fmt.Println(fileBackups)
	// # Rolling
	// # 1: find first empty cell
	j := B
	i := 1
	for i = 1; i < B; i++ {
		if fileBackups[i] == "" {
			j = i
			break
		}
	}
	// I have all files
	if j == B {
		// Delete last file
		err := os.Remove(filepath.Join(path, fileBackups[j-1]))
		if err != nil {
			WriteError(err)
		}
		j--
	}
	// Now I can rotate all file
	for i = j - 1; i > 0; i-- {
		oldPath := filepath.Join(path, fileBackups[i])
		newPath := filepath.Join(path, comparator.Replace(j))
		if err := os.Rename(oldPath, newPath); err != nil {
			WriteError(err)
		}
		j = i
	}
	// Now I can Rotate the hot file
	oldPath := filepath.Join(path, hotFilename)
	newPath := filepath.Join(path, comparator.Replace(1))
	if err := os.Rename(oldPath, newPath); err != nil {
		WriteError(err)
	}
}

const ChannelSize = 160

func FileRotationIV(filename string, backups uint) {
	// Opening folder
	path, hotFilename, _ := getPathHotFileExt(filename)
	comp := NewComparator(hotFilename, ".", false, false)
	backups++

	filesChannel, _ := readFiles(path)
	chValidFiles, chDeleteFiles := matchOrDelete(filesChannel, comp, backups)
	done := deleteFiles(chDeleteFiles, path)

	presentFiles := make([]bool, backups)
	files := make([]string, backups)
	firstRenameableFile := 1
	for f := range chValidFiles {
		presentFiles[f.value] = true
		if f.value == firstRenameableFile {
			firstRenameableFile++
			for firstRenameableFile < int(backups) && presentFiles[firstRenameableFile] {
				firstRenameableFile++
			}
		}
		files[f.value] = f.name
	}
	if firstRenameableFile == int(backups) {
		if err := os.Remove(filepath.Join(path, files[firstRenameableFile-1])); err != nil {
			WriteError(err)
		}
		firstRenameableFile--
	}
	firstRenameableFile--

	next := filepath.Join(path, files[firstRenameableFile])
	for firstRenameableFile = firstRenameableFile - 1; firstRenameableFile > 0; firstRenameableFile-- {
		curr := filepath.Join(path, comp.Replace(firstRenameableFile))
		if err := os.Rename(next, curr); err != nil {
			WriteError(err)
		}
		next = curr
	}
	// Now I can Rotate the hot file
	if err := os.Rename(filepath.Join(path, hotFilename), next); err != nil {
		WriteError(err)
	}
	<-done
}

func readFiles(path string) (chan string, error) {
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

func matchOrDelete(files chan string, comp *Comparator, numBackups uint) (chan struct {
	value int
	name  string
}, chan string) {
	chInt := make(chan struct {
		value int
		name  string
	}, ChannelSize)
	chStr := make(chan string, ChannelSize)
	go func() {
		for file := range files {
			m := comp.Match(file)
			if m > -1 {
				if m < int(numBackups) {
					chInt <- struct {
						value int
						name  string
					}{value: m, name: file}
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

func deleteFiles(chFiles chan string, path string) chan interface{} {
	done := make(chan interface{})
	go func() {
		for file := range chFiles {
			if err := os.Remove(filepath.Join(path, file)); err != nil {
				WriteError(err)
			}
		}
		close(done)
	}()
	return done
}

// FileRotationV fifth algorithm
func FileRotationV(filename string, backups uint) {
	// Opening folder
	path, hotFilename, _ := getPathHotFileExt(filename)
	comp := NewComparator(hotFilename, ".", false, false)

	// read files
	filesChannel, _ := readFiles(path)
	// match file or delete it
	chValidFiles, chDeleteFiles := matchOrDeleteV(filesChannel, comp, int(backups))

	// delete files in a goroutine (WARNING)
	done := deleteFiles(chDeleteFiles, path)

	// presents files
	presentFiles := make([]bool, backups)
	for f := range chValidFiles {
		presentFiles[f] = true
	}

	// get first renameable file
	firstRenameableFile := 1
	for firstRenameableFile < len(presentFiles) {
		if !presentFiles[firstRenameableFile] {
			break
		}
		firstRenameableFile++
	}

	// change current working directory
	wd, _ := os.Getwd()
	if err := os.Chdir(path); err != nil {
		WriteError(err)
		return
	}

	// rename from the last filename
	for firstRenameableFile--; firstRenameableFile > 0; firstRenameableFile-- {
		if err := os.Rename(comp.Replace(firstRenameableFile), comp.Replace(firstRenameableFile+1)); err != nil {
			WriteError(err)
		}
	}

	// now I can Rotate the hot file
	if err := os.Rename(hotFilename, comp.Replace(1)); err != nil {
		WriteError(err)
	}

	// set last working directory
	if err := os.Chdir(wd); err != nil {
		WriteError(err)
		return
	}

	<-done
}

func matchOrDeleteV(files chan string, comp *Comparator, numBackups int) (chan int, chan string) {
	chInt := make(chan int, ChannelSize)
	chStr := make(chan string, ChannelSize)
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

// FileRotationVI sixth algorithm
func FileRotationVI(filename string, backups uint) {
	// Opening folder
	path, hotFilename, _ := getPathHotFileExt(filename)
	comp := NewComparator(hotFilename, ".", false, false)

	// read files
	filesChannel, _ := readFiles(path)
	// match file or delete it
	chValidFiles, chDeleteFiles := matchOrDeleteV(filesChannel, comp, int(backups))

	// delete files in a goroutine (WARNING)
	done := deleteFiles(chDeleteFiles, path)

	// presents files
	presentFiles := make([]bool, backups)
	for f := range chValidFiles {
		presentFiles[f] = true
	}

	// get first renameable file
	firstRenameableFile := 1
	for firstRenameableFile < len(presentFiles) {
		if !presentFiles[firstRenameableFile] {
			break
		}
		firstRenameableFile++
	}
	next := path + comp.Replace(firstRenameableFile-1)
	for firstRenameableFile = firstRenameableFile - 2; firstRenameableFile > 0; firstRenameableFile-- {
		curr := path + comp.Replace(firstRenameableFile)
		if err := os.Rename(next, curr); err != nil {
			WriteError(err)
		}
		next = curr
	}
	// Now I can Rotate the hot file
	if err := os.Rename(path+hotFilename, next); err != nil {
		WriteError(err)
	}
	<-done
}
