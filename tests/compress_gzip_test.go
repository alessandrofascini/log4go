package tests

import (
	"compress/gzip"
	"fmt"
	"os"
	"testing"
)

func TestCompress(*testing.T) {
	// Some text we want to compress.
	original := "bird and frog"

	// Open a file for writing.
	f, _ := os.Create("./tests/file.log")

	// Create gzip writer.
	w := gzip.NewWriter(f)

	// Write bytes in compressed form to the file.
	w.Write([]byte(original))

	// Close the file.
	w.Close()

	fmt.Println("DONE")
}
