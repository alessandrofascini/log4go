package appenders

import (
	"fmt"
	"sync"

	"github.com/alessandrofascini/log4go/pkg"
)

type AppenderPool map[string]*AppenderHandler

// NewAppenderPool This function create a pool of appenders. Returns a map[string]Appender.
// In other words, this function create a concrete map starting from a configuration
func NewAppenderPool(c pkg.AppendersConfig, globalAppenderWaitGroup *sync.WaitGroup) AppenderPool {
	a := AppenderPool{}
	for key, config := range c {
		if a[key] != nil { // already created
			panic(fmt.Sprintf("already exists appender with name %q", key))
		}
		a[key] = NewAppenderHandler(config, globalAppenderWaitGroup)
	}
	return a
}
