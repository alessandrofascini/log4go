package layouts_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/alessandrofascini/log4go/v2/events"
	"github.com/alessandrofascini/log4go/v2/layouts"
	"github.com/alessandrofascini/log4go/v2/levels"
)

func TestNewLayout(t *testing.T) {
	l, err := layouts.NewLayout("%h@[%d] - %m{1,5} %m %m\n")
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}

	var buf bytes.Buffer
	l.WriteEventTo(&buf, &events.Event{
		StartTime: time.Now(),
		Data:      []any{map[string]any{"ciao": 132}, "hello1", "hello2"},
		Level:     *levels.LevelInfo,
	})
	fmt.Println(buf.String())
}

func BenchmarkNewLayout(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := layouts.NewLayout("%h@[%d] - %m{1,5} %m %m\n")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkWriteEventTo(b *testing.B) {
	l, err := layouts.NewLayout("%h@[%d] - %m{1,5} %m %m\n")
	if err != nil {
		b.Error(err)
		b.Fail()
	}
	var buf bytes.Buffer

	data := []any{map[string]any{"ciao": 132}, "hello1", "hello2"}

	for i := 0; i < b.N; i++ {
		l.WriteEventTo(&buf, &events.Event{
			StartTime: time.Now(),
			Data:      data,
			Level:     *levels.LevelInfo,
		})
	}

}
