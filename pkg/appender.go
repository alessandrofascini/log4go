package pkg

import (
	"fmt"
	"strconv"

	configuration "github.com/alessandrofascini/log4go/internal/config"
)

type Appender func(LoggingEvent)

type AppenderConfig map[string]any

func (c AppenderConfig) GetLayout() LayoutConfig {
	switch v := c["layout"].(type) {
	case nil:
		return LayoutConfig{
			"type": "coloured",
		}
	case LayoutConfig:
		return v
	}
	panic("layout miss type check")
}

func (c AppenderConfig) GetType() string {
	v, ok := c["type"]
	if !ok {
		panic(fmt.Errorf("cannot read property %q in %v", "type", v))
	}
	return fmt.Sprintf("%v", v)
}

func (c AppenderConfig) GetChannelSize() int {
	v, ok := c["channelSize"]
	if !ok {
		return configuration.GetDefaultAppenderChannelSize()
	}
	switch t := v.(type) {
	case int:
		if t > 0 {
			return t
		}
	case string:
		if value, err := strconv.Atoi(t); err != nil {
			return value
		}
	}
	return configuration.GetDefaultAppenderChannelSize()
}

func (c AppenderConfig) GetInternalChannelSize() int {
	v, ok := c["channelSize"]
	if !ok {
		return configuration.GetDefaultInternalChannelSize()
	}
	switch t := v.(type) {
	case int:
		if t > 0 {
			return t
		}
	case string:
		if value, err := strconv.Atoi(t); err != nil {
			return value
		}
	}
	return configuration.GetDefaultInternalChannelSize()
}

type AppendersConfig map[string]AppenderConfig
