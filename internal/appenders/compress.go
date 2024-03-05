package appenders

import (
	"compress/gzip"
	"os"
)

/*
this function compress file.
old is the old filename
curr is the new filename after compression
level is the compress mode
*/
func compress(old, curr string, compressMode int) error {
	// create new file
	f, err := os.Create(curr)
	if err != nil {
		return err
	}
	// create new writer from new file
	w, err := gzip.NewWriterLevel(f, compressMode)
	if err != nil {
		return err
	}
	// Read file content
	content, err := os.ReadFile(old)
	if err != nil {
		return err
	}
	// Write all content
	if _, err = w.Write(content); err != nil {
		return err
	}
	// Close the file
	if err = w.Close(); err != nil {
		return err
	}
	// remove old file
	if err = os.Remove(old); err != nil {
		return err
	}
	return nil
}
