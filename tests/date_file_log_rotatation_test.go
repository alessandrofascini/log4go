package tests

import (
	"fmt"
	"os"
	"testing"
	"time"
)

const dateBase = base

func generateDateFiles(base string, n int) {
	path, _, _ := getPathHotFileExt(base)
	os.Create(path)
	os.Create(base)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	for i := 0; i < n; i++ {
		today = today.Add(-1 * time.Hour * 24)
		filename := fmt.Sprintf("%s.%s", base, today.Format("2006-01-02"))
		if _, err := os.Create(filename); err != nil {
			WriteError(err)
		}
	}
}

func TestGenerateDateFiles(t *testing.T) {
	generateDateFiles(dateBase, 10)
}

func generateDateFilesWithIndex(base string, n int) {
	path, _, _ := getPathHotFileExt(base)
	os.Mkdir(path, 0777)
	os.Create(base)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	for i := 0; i < n; i++ {
		for j := 1; j < n; j++ {
			filename := fmt.Sprintf("%s.%s.%d", base, today.Format("2006-01-02"), j)
			if _, err := os.Create(filename); err != nil {
				WriteError(err)
			}
		}
		today = today.Add(-1 * time.Hour * 24)
	}
}

func TestGenerateDateFilesWithIndex(t *testing.T) {
	generateDateFilesWithIndex(dateBase, 10)
}

func TestDateFileRotationIICaseI(t *testing.T) {
	generateDateFiles(dateBase, 10)
	DateFileRotationIICaseI(dateBase, 5)
}

func TestDateFileRotationIICaseII(t *testing.T) {
	generateDateFilesWithIndex(dateBase, 15)
	generateFiles(dateBase, 10)
	DateFileRotationIICaseII(dateBase, 20)
}
