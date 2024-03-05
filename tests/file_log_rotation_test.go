package tests

import (
	"fmt"
	"os"
	"testing"
)

const N = 1000
const base = "./tests/temp/mon.log"

func TestGenerateFiles(t *testing.T) {
	generateFiles(base, N)
}

func generateFiles(base string, n int) {
	os.Create(base)
	for i := 1; i <= n; i++ {
		os.Create(fmt.Sprintf("%s.%d", base, i))
	}
}

// III

func TestFileRotationIII(t *testing.T) {
	FileRotationIII(base, N)
}

func BenchmarkFileRotationIII(b *testing.B) {
	generateFiles(base, N)
	b.StartTimer()
	FileRotationIII(base, N)
	b.StopTimer()
}

// IV

func TestFileRotationIV(t *testing.T) {
	FileRotationIV(base, N)
}

func BenchmarkFileRotationIV(b *testing.B) {
	generateFiles(base, N)
	b.StartTimer()
	FileRotationIV(base, N)
	b.StopTimer()
}

// V

func TestFileRotationV(t *testing.T) {
	generateFiles(base, N)
	FileRotationV(base, N)
}

func BenchmarkFileRotationV(b *testing.B) {
	generateFiles(base, N)
	b.StartTimer()
	FileRotationV(base, N)
	b.StopTimer()
}

// VI

func TestFileRotationVI(t *testing.T) {
	generateFiles(base, N)
	FileRotationVI(base, N)
}

func BenchmarkFileRotationVI(b *testing.B) {
	generateFiles(base, N)
	b.StartTimer()
	FileRotationVI(base, N)
	b.StopTimer()
}

// Compare All

func BenchmarkFunctions(b *testing.B) {
	b.Run("third", BenchmarkFileRotationIII)
	b.Run("fourth", BenchmarkFileRotationIV)
	b.Run("fifth", BenchmarkFileRotationV)
	b.Run("sixth", BenchmarkFileRotationVI)
}
