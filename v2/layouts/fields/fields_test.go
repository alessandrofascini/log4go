package fields_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/alessandrofascini/log4go/v2/events"
	"github.com/alessandrofascini/log4go/v2/layouts/fields"
	"github.com/alessandrofascini/log4go/v2/levels"
)

// TODO: Complete me

func TestNewField(t *testing.T) {
	testcases := []string{"m", "5m", ".15m", ".15", "-15m", "-20m", "15.-26m", "-9.31m", "m{1,5}", "[", "m", "]", "r", "h", "d"}

	var buf bytes.Buffer
	buf.WriteByte('\n')

	for i := range testcases {
		t.Log(testcases[i])
		reader := strings.NewReader(testcases[i])
		field, err := fields.NewField(reader)
		t.Log(field, err)
		if err == nil {
			eReader, _ := fields.NewEventReader(&events.Event{
				Data:      []any{"ciasdsadkkasl"},
				StartTime: time.Now(),
				Level:     *levels.LevelInfo,
			})
			buf.WriteString(fmt.Sprintf("%10s: ", testcases[i]))
			field.WriteTo(&buf, eReader)
			buf.WriteByte('\n')
		}
	}

	t.Log(buf.String())
}
