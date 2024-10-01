package appenders

import (
	"context"
	"github.com/alessandrofascini/log4go/v2/events"
	"github.com/alessandrofascini/log4go/v2/layouts"
	"github.com/alessandrofascini/log4go/v2/levels"
	"os"
	"testing"
	"time"
)

func TestUnchanneled_WriteEvent(t *testing.T) {
	layout, err := layouts.NewLayout("ciao %d%n")
	if err != nil {
		t.Fatal(err)
	}
	app, err := NewUnchanneled(os.Stdout, layout)
	if err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < 1000; i++ {
		_, err := app.WriteEvent(&events.Event{
			StartTime:    time.Now(),
			CategoryName: "test",
			Data:         make([]any, 0),
			Level:        levels.LevelTrace,
			Context:      context.TODO(),
			Pid:          15323,
		})
		if err != nil {
			t.Error(err)
		}
	}

	if err = app.Close(); err != nil {
		t.Error(err)
	}
}
