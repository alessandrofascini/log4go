package v2_test

import (
	"fmt"
	"io"
	"os"
	"testing"
)

type Writer struct {
	*os.File
}

func NewWriter() *Writer {
	file, _ := os.Open("./bench.txt")
	return &Writer{file}
}

var writer io.Writer = NewWriter()

type LevelWriter interface {
	Trace(...any)
	Debug(...any)
	Info(...any)
	Warn(...any)
	Error(...any)
	Fatal(...any)
}

// implements all methods
type TraceWriter struct{}

func (w *TraceWriter) log(_ int, args []any) {
	fmt.Fprint(writer, args...)
}

func (w *TraceWriter) Trace(args ...any) {
	w.log(0, args)
}

func (w *TraceWriter) Debug(args ...any) {
	w.log(1, args)
}

func (w *TraceWriter) Info(args ...any) {
	w.log(2, args)
}

func (w *TraceWriter) Warn(args ...any) {
	w.log(3, args)
}

func (w *TraceWriter) Error(args ...any) {
	w.log(4, args)
}

func (w *TraceWriter) Fatal(args ...any) {
	w.log(5, args)
}

type DebugWriter struct {
	TraceWriter
}

func (d *DebugWriter) Trace(args ...any) {
	// do nothing
}

type Logger struct {
	minLevel int
}

func (l *Logger) log(_ int, args []any) {
	fmt.Fprint(writer, args)
}

func (w *Logger) Trace(args ...any) {
	if w.minLevel < 0 {
		return
	}
	w.log(0, args)
}

func (w *Logger) Debug(args ...any) {
	if w.minLevel < 1 {
		return
	}
	w.log(1, args)
}

func (w *Logger) Info(args ...any) {
	if w.minLevel < 2 {
		return
	}
	w.log(2, args)
}

func (w *Logger) Warn(args ...any) {
	if w.minLevel < 3 {
		return
	}
	w.log(3, args)
}

func (w *Logger) Error(args ...any) {
	if w.minLevel < 4 {
		return
	}
	w.log(4, args)
}

func (w *Logger) Fatal(args ...any) {
	if w.minLevel < 5 {
		return
	}
	w.log(5, args)
}

func BenchmarkLog(b *testing.B) {
	tw, dw := new(TraceWriter), new(DebugWriter)

	b.Run("wihout if", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tw.Trace("trace", "trace")
			dw.Trace("trace", "debug")
			tw.Debug("debug", "trace")
			dw.Debug("debug", "debug")
		}
	})

	tw1, dw1 := &Logger{0}, &Logger{1}

	b.Run("with if", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tw1.Trace("trace", "trace")
			dw1.Trace("trace", "debug")
			tw1.Debug("debug", "trace")
			dw1.Debug("debug", "debug")
		}
	})
}
