package tests

import (
	"fmt"
	"os"
	"testing"
)

const thisFolder = ".log4Go/tests"

func TestFlags(t *testing.T) {
	fmt.Printf("read only %b\n", os.O_RDONLY)
	fmt.Printf("write only %b\n", os.O_WRONLY)
	fmt.Printf("read and write %b\n", os.O_RDWR)
	fmt.Printf("append %b\n", os.O_APPEND)
	fmt.Printf("create %b\n", os.O_CREATE)
	fmt.Printf("exec %b\n", os.O_EXCL)
	fmt.Printf("sync %b\n", os.O_SYNC)
	fmt.Printf("trun %b\n", os.O_TRUNC)
}

// a+
func TestReadAndAppend(t *testing.T) {
	const flag = os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile("readAndAppend.txt", flag, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	if _, err := file.WriteString("Hello from Code!"); err != nil {
		panic(err)
	}
}
