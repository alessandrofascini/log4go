package appenders

import (
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"github.com/alessandrofascini/log4go/pkg"

	errorshelper "github.com/alessandrofascini/log4go/internal/errors"
)

type FileManagerFactory func(config *FileConfig) IFileManagerRotation

type FileConfig struct {
	filename     string
	flags        int
	mode         os.FileMode
	maxLogSize   int64
	backups      int
	compress     bool
	compressMode int
	keepFileExt  bool
	fileNameSep  string
	pattern      string
	channelSize  int
}

func newFileConfig(fc pkg.AppenderConfig) *FileConfig {
	config := &FileConfig{
		flags:        os.O_APPEND | os.O_CREATE | os.O_WRONLY,
		mode:         0755,
		maxLogSize:   0,
		backups:      5,
		compress:     false,
		compressMode: gzip.DefaultCompression,
		keepFileExt:  false,
		fileNameSep:  ".",
		pattern:      "2006-01-02",
		channelSize:  fc.GetInternalChannelSize(),
	}

	// filename
	// get filename and validate it
	func(fc pkg.AppenderConfig, conf *FileConfig) {
		var filename string
		switch v := fc["filename"].(type) {
		case string:
			filename, _ = filepath.Abs(filepath.Clean(v))
		default:
			panic(fmt.Errorf("missing or invalid %q attribute for appender config", "filename"))
		}
		dir, file := filepath.Split(filename)
		if dir == "" {
			dir, _ = os.Getwd()
		}
		ext := filepath.Ext(file)
		if ext == "" {
			file = file + ".log"
		}
		if ext == ".gz" {
			panic(fmt.Errorf("cannot use a file with extension %q. Please choose another file extension", ".gz"))
		}
		config.filename = filepath.Join(dir, file)
	}(fc, config)

	// max log size
	switch v := fc["maxLogSize"].(type) {
	case int:
		if v < 0 {
			panic("maxLogSize must be greater than zero")
		}
		config.maxLogSize = int64(v)
	case float64:
		if v < 0 {
			panic("maxLogSize must be greater than zero")
		}
		config.maxLogSize = int64(v)
	case string:
		values := []int{
			'M': 20,
			'K': 10,
			'G': 30,
		}
		m := len(v) - 1
		unit := v[m]
		if unicode.IsDigit(rune(unit)) {
			errorshelper.WriteErrorf("maxLogSize: %q is invalid", v)
			break
		}
		numberStr := v[0:m]
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			errorshelper.WriteErrorf("maxLogSize: %q is invalid", v)
			break
		}
		config.maxLogSize = int64(number) << values[unit]
	default:
		panic("invalid type of maxLogSize")
	}

	// backups
	switch v := fc["backups"].(type) {
	case int:
		if v < 0 {
			panic("backups must be greater than zero")
		}
		config.backups = v
	}

	// compress
	switch v := fc["compress"].(type) {
	case bool:
		config.compress = v
	}

	// compress mode
	switch v := fc["compressMode"].(type) {
	case int:
		switch v {
		case gzip.DefaultCompression, gzip.BestSpeed, gzip.BestCompression, gzip.HuffmanOnly, gzip.NoCompression:
			config.compressMode = v
		}
	case string:
		switch strings.ToLower(v) {
		default:
			fallthrough
		case "default", "defaultcompression":
			config.compressMode = gzip.DefaultCompression
		case "bestspeed":
			config.compressMode = gzip.BestSpeed
		case "bestcompression":
			config.compressMode = gzip.BestCompression
		case "huffmanonly":
			config.compressMode = gzip.HuffmanOnly
		case "nocompression":
			config.compressMode = gzip.NoCompression
		}
	}

	// keep file ext
	switch v := fc["keepFileExt"].(type) {
	case bool:
		config.keepFileExt = v
	default:
		panic("invalid type of compress")
	}

	switch v := fc["channelSize"].(type) {
	case int:
		if v >= 0 {
			config.channelSize = v
		}
	case string:
		value, err := strconv.Atoi(v)
		if err != nil {
			break
		}
		if value >= 0 {
			config.channelSize = value
		}
	}

	switch v := fc["fileNameSep"].(type) {
	case string:
		config.fileNameSep = v
	}

	// EXPERIMENTAL
	switch v := fc["mode"].(type) {
	case int:
		config.mode = os.FileMode(v)
	}

	// EXPERIMENTAL
	switch v := fc["flags"].(type) {
	case int:
		config.flags = v
	}

	// TODO Pattern (future feature)

	return config
}
