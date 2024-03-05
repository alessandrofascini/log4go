package tests

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestOpenFile(t *testing.T) {
	const filename = "./tests/readme.txt"
	file, _ := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0600)
	if _, err := file.WriteString("\nanother string"); err != nil {
		panic(err)
	}
	info, _ := file.Stat()
	fmt.Println("after", info.Size(), info.ModTime())
}

func TestFileReadDir(t *testing.T) {
	const path = "./tests/temp"
	const filename = "mon.log"
	//file, _ := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_RDONLY, 0644)
	file, _ := os.Open(path)
	fmt.Println(file.Name())
	names, err := file.Readdirnames(0)
	if err != nil {
		fmt.Println(fmt.Errorf("%+v", err))
		return
	}

	// Sort folder
	sort.Strings(names)

	// Reversed It
	left := 0
	right := len(names) - 1
	for left < right {
		t := names[left]
		names[left] = names[right]
		names[right] = t
		left++
		right--
	}

	// Rename not hot file
	for i, name := range names {
		if v, _ := filepath.Match(fmt.Sprintf("%s.*", filename), name); v {
			fmt.Println(name, "match!")
			//Pick the number
			counter := len(names) - i
			//Change Name
			fmt.Println(filepath.Join(path, name), filepath.Join(path, fmt.Sprintf("%s.%d", filename, counter)))
			if err := os.Rename(filepath.Join(path, name), filepath.Join(path, fmt.Sprintf("%s.%d", filename, counter))); err != nil {
				panic(err)
			}
		}
	}

	// renaming the hot file
	if err := os.Rename(filepath.Join(path, filename), filepath.Join(path, fmt.Sprintf("%s.%d", filename, 1))); err != nil {
		panic(err)
	}
}

func TestCompressFiles(t *testing.T) {
	// Open file on disk.
	const path = "./tests/temp"
	filename := "mon.log.1"
	f, _ := os.Open(filepath.Join(path, filename))

	// Create a Reader and use ReadAll to get all the bytes from the file.
	reader := bufio.NewReader(f)
	content, _ := io.ReadAll(reader)

	// Replace txt extension with gz extension.
	filename = fmt.Sprintf("%s.gz", filename)

	// Open file for writing.
	f, _ = os.Create(filepath.Join(path, filename))

	// Write compressed data.
	w := gzip.NewWriter(f)
	w.Write(content)
	w.Close()

	// TODO delete old file

	// Done.
	fmt.Println("DONE")
}

func TestRenameFile(t *testing.T) {
	const path = "./tests/temp"
	if err := os.Rename(filepath.Join(path, "mon.log.1"), filepath.Join(path, "mon.log.1")); err != nil {
		panic(err)
	}
	fmt.Println("renamed!")
}
