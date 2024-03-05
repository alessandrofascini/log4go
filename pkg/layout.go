package pkg

import (
	"fmt"
)

type Layout func(LoggingEvent) string
type LayoutConfig map[string]any

func (c LayoutConfig) GetType() string {
	v := c["type"]
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

type LayoutTokens map[string]any
