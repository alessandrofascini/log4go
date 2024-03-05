package main

import (
	"fmt"

	"github.com/alessandrofascini/log4go/internal/layouts"
	log "github.com/alessandrofascini/log4go/internal/types"

	"github.com/alessandrofascini/log4go/internal/layouts/types"
)

type ConcreteLayoutConfiguration struct {
}

func (c ConcreteLayoutConfiguration) GetType() types.LayoutType {
	return types.LayoutTypePattern
}

func (c ConcreteLayoutConfiguration) GetPattern() types.PatternConfig {
	return "%m %x{user} %d %d{ISO8601_WITH_TZ_OFFSET} %d{ABSOLUTETIME} %d{DATETIME}"
}

type ConcreteContext struct {
}

func (c ConcreteContext) GetLevel() string {
	return log.FATAL.String()
}

func (c ConcreteContext) GetCategoryName() string {
	return "default"
}

func (c ConcreteContext) GetLogData() string {
	return "Hello From Alessandro"
}

func (c ConcreteContext) GetTokens() types.Tokens {
	return types.Tokens{
		"user": func() any {
			return "This is a function user"
		},
	}
}

func main() {
	//fmt.Printf("%v %v\n", "Alessandro", 2003)
	//pattern := "[%d] [%p] %c - %m"
	//result := layouts_old.CreateLayout(pattern)
	//fmt.Println(result.GetFormat())
	//
	//logContext := entities.LogEntity{}
	//logContext.SetLevel(entities.INFO)
	//logContext.SetMessage("Hello From Alessandro")
	//logContext.SetCategoryName("default")
	//
	//fmt.Println(result.GetPattern(&layouts_old.EnvContext{
	//	LogContext: logContext,
	//}))
	conf := &ConcreteLayoutConfiguration{}
	layout := layouts.LayoutFactory(conf)
	fmt.Println(layout.Fill(&ConcreteContext{}))
}
