package fields

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrUnkownVerb  = errors.New("unkown field verb")
	ErrTruncFormat = errors.New("incorrect truncation format")
	ErrPadFormat   = errors.New("incorrect padding format")
)

type Field interface {
	WriteTo(io.Writer, *EventReader) (n int, err error)
}

const (
	// %r time in toLocaleTimeString format
	localeTimeVerb rune = 'r'
	// %p log level
	logLevelVerb rune = 'p'
	// %c log category
	logCategoryVerb rune = 'c'
	// %h hostname
	hostnameVerb rune = 'h'
	// %m log data
	logDataVerb rune = 'm'
	// %d date, formatted - default is ISO8601, format options are: ISO8601, ISO8601_WITH_TZ_OFFSET, ABSOLUTETIME, DATETIME, or any string compatible with the date-format library. e.g. %d{DATETIME}, %d{yyyy/MM/dd-hh.mm.ss}
	dateVerb rune = 'd'
	// %% % - for when you want a literal % in your output
	percentVerb rune = '%'
	// %n newline
	newlineVerb rune = 'n'
	// %z process id (from process.pid)
	processIdVerb rune = 'z'
	// %[ start a coloured block (colour will be taken from the log level, similar to colouredLayout)
	startColouredBlockVerb rune = '['
	// %] end a coloured block
	endColouredBlockVerb rune = ']'
)

var verbTable = map[rune]struct{}{
	localeTimeVerb:         {},
	logLevelVerb:           {},
	logCategoryVerb:        {},
	hostnameVerb:           {},
	logDataVerb:            {},
	dateVerb:               {},
	percentVerb:            {},
	newlineVerb:            {},
	processIdVerb:          {},
	startColouredBlockVerb: {},
	endColouredBlockVerb:   {},
}

func isValidVerb(r rune) bool {
	_, ok := verbTable[r]
	return ok
}

func NewField(r *strings.Reader) (f Field, err error) {
	pad, padOk, err := readPadding(r)
	if err != nil {
		return nil, err
	}

	trunc, truncOk, err := readTruncation(r)
	if err != nil {
		return nil, err
	}

	verb, err := readVerb(r)
	if err != nil {
		return nil, err
	}

	format, err := readFormat(r)
	if err != nil {
		return nil, err
	}

	switch verb {
	case localeTimeVerb:
		f = NewLocaleTimeField()
	case logLevelVerb:
		f = &LogLevelField{}
	case logCategoryVerb:
		f = &LogCategoryField{}
	case logDataVerb:
		f = NewMessageField(format)
	case hostnameVerb:
		f = &HostnameField{}
	case dateVerb:
		f = NewDateField(format)
	case percentVerb:
		f = NewPercentField()
	case newlineVerb:
		f = NewNewLineField()
	case processIdVerb:
		f = &ProcessIdField{}
	case startColouredBlockVerb:
		f = &ColorField{}
	case endColouredBlockVerb:
		f = &ColorResetField{}
	default:
		return nil, ErrUnkownVerb
	}

	if padOk {
		f = NewPaddingField(f, pad)
	}

	if truncOk {
		f = &TruncationField{f, trunc}
	}

	return
}

func readPadding(r *strings.Reader) (int, bool, error) {
	c, _, err := r.ReadRune()
	if err != nil {
		if err == io.EOF {
			r.UnreadRune()
			return 0, false, nil
		}
		return 0, false, err
	}

	number := make([]rune, 0, 10)

	switch {
	default:
		fallthrough
	case c == '.', isValidVerb(c):
		r.UnreadRune()
		return 0, false, nil
	case c == '-', unicode.IsDigit(c):
		number = append(number, c)
	}

	for r.Len() > 0 {
		c, _, err = r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, false, err
		}
		if unicode.IsDigit(c) {
			number = append(number, c)
			continue
		}
		if c == '.' || isValidVerb(c) {
			r.UnreadRune()
			break
		}
		return 0, false, ErrPadFormat
	}

	n, err := strconv.Atoi(string(number))
	if err != nil {
		return 0, false, nil
	}
	return n, false, nil
}

func readTruncation(r *strings.Reader) (int, bool, error) {
	c, _, err := r.ReadRune()
	if err != nil {
		if err == io.EOF {
			r.UnreadRune()
			return 0, false, nil
		}
		return 0, false, err
	}
	switch {
	case isValidVerb(c):
		r.UnreadRune()
		return 0, false, nil
	case c != '.':
		return 0, false, ErrTruncFormat
	}

	c, _, err = r.ReadRune()
	if err != nil {
		return 0, false, err
	}
	number := make([]rune, 0, 10)
	switch {
	case c == '+':
	case c == '-', unicode.IsDigit(c):
		number = append(number, c)
	default:
		r.UnreadRune()
		return 0, false, ErrTruncFormat
	}
	for r.Len() > 0 {
		c, _, err = r.ReadRune()
		if err != nil {
			return 0, false, err
		}
		if unicode.IsDigit(c) {
			number = append(number, c)
			continue
		}
		if isValidVerb(c) {
			r.UnreadRune()
			break
		}
		return 0, false, ErrTruncFormat
	}
	n, err := strconv.Atoi(string(number))
	if err != nil {
		return 0, false, err
	}
	return n, false, nil
}

func readVerb(r *strings.Reader) (rune, error) {
	c, _, err := r.ReadRune()
	if err != nil {
		return c, err
	}
	if !isValidVerb(c) {
		return c, ErrUnkownVerb
	}
	return c, nil
}

func readFormat(r *strings.Reader) (string, error) {
	c, _, _ := r.ReadRune()
	if c != '{' {
		r.UnreadRune()
		return "", nil
	}
	format := make([]rune, 0, 15)
	for r.Len() > 0 {
		c, _, err := r.ReadRune()
		if err != nil {
			return "", err
		}
		if c == '}' {
			break
		}
		format = append(format, c)
	}
	return string(format), nil
}
