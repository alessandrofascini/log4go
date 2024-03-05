package categories

import (
	"fmt"

	"github.com/alessandrofascini/log4go/pkg"

	"github.com/alessandrofascini/log4go/internal/appenders"
)

type Category struct {
	CategoryName string
	Appenders    []*appenders.AppenderHandler
	LevelValue   uint
}

// NewCategory This function create a Category, consuming a Category and a pool of Appender.
// Returns: a Category
func NewCategory(config pkg.CategoryConfig, pool appenders.AppenderPool) Category {
	// Create an array of Appenders
	var a []*appenders.AppenderHandler
	names := config.GetAssociatedAppenders()
	for _, name := range names {
		if pool[name] == nil {
			panic(fmt.Sprintf("not found appender with name %q", name))
		}
		a = append(a, pool[name])
	}
	return Category{
		CategoryName: config.GetName(),
		Appenders:    a,
		LevelValue:   config.GetLevelValue(),
	}
}
