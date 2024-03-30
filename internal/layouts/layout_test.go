package layouts_test

import (
	"testing"
	"time"

	"github.com/alessandrofascini/log4go/internal/layouts"
	"github.com/alessandrofascini/log4go/pkg"
)

func BenchmarkNewLayout(b *testing.B) {
	c := map[string]any{}

	c["pattern"] = "%h@[%d] - %m{1,5} %m %m\n"

	config := pkg.LayoutConfig(c)

	for i := 0; i < b.N; i++ {
		layouts.PatternLayout(&config)
	}
}

func BenchmarkWriteEvent(b *testing.B) {
	c := map[string]any{}

	c["pattern"] = "%h@[%d] - %m{1,5} %m %m\n"

	config := pkg.LayoutConfig(c)

	l := layouts.PatternLayout(&config)

	data := []any{map[string]any{"ciao": 132}, "hello1", "hello2"}

	for i := 0; i < b.N; i++ {
		l(pkg.LoggingEvent{
			StartTime:    time.Now(),
			Data:         data,
			CategoryName: "",
			Level:        *pkg.LevelInfo,
			Context:      nil,
			Pid:          12,
		})
	}

}
