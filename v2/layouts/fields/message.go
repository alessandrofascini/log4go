package fields

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"
)

type MessageField struct {
	left, right int
}

func NewMessageField(limits string) *MessageField {
	l, r := limits, ""
	if i := strings.IndexRune(limits, ','); i != -1 {
		l, r = limits[:i], limits[i+1:]
	}
	left, err := strconv.Atoi(string(l))
	if err != nil {
		left = 0
	}
	right, err := strconv.Atoi(string(r))
	if err != nil {
		right = -1
	}
	return &MessageField{left, right}
}

func (m *MessageField) runesToUTF8Manual(rs []rune) []byte {
	size := 0
	for _, r := range rs {
		size += utf8.RuneLen(r)
	}

	bs := make([]byte, size)

	count := 0
	for _, r := range rs {
		count += utf8.EncodeRune(bs[count:], r)
	}

	return bs
}

func (m *MessageField) WriteTo(w io.Writer, reader *EventReader) (n int, err error) {
	data, err := reader.ReadData()
	if err != nil {
		if err == ErrEmptyReader {
			return w.Write([]byte{})
		}
		return 0, err
	}

	b := []rune(fmt.Sprintf("%v", data))
	l := len(b)
	right := m.right
	if m.right == -1 {
		right = l
	}
	left := m.left
	if left > l {
		left = l
	}
	return w.Write(m.runesToUTF8Manual(b[left:right]))
}
