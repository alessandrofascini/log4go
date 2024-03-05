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
