package log4go

import (
	"fmt"

	"github.com/alessandrofascini/log4go/pkg"
)

func validateLayoutConfig(c map[string]any) pkg.LayoutConfig {
	// parameter required: type
	switch c["type"].(type) {
	case string:
		return c
	}
	panic("invalid layout config")
}

func validateAppendersConfig(c map[string]any) pkg.AppendersConfig {
	res := pkg.AppendersConfig{}
	set := map[string]bool{}
	for appenderName, config := range c {
		// unique appender name
		if set[appenderName] {
			panic("there are two or more appenders with the same name!")
		}
		set[appenderName] = true
		switch v := config.(type) {
		case map[string]any:
			res[appenderName] = validateAppenderConfig(v)
		default:
			panic(fmt.Sprintf("invalid type for appender with name %q: required map[string]any, found %T", appenderName, v))
		}
	}
	return res
}

func validateAppenderConfig(c map[string]any) pkg.AppenderConfig {
	switch c["type"].(type) {
	case string:
	default:
		panic(fmt.Sprintf("missing required field %q from appender configuration", "type"))
	}
	switch v := c["layout"].(type) {
	case nil:
		// default settings
		c["layout"] = pkg.LayoutConfig{
			"type": "basic",
		}
	case map[string]any:
		c["layout"] = validateLayoutConfig(v)
	default:
		panic(fmt.Sprintf("layout type field provided is incorrect! required string, found %T", v))
	}
	return c
}

func validateCategoriesConfig(c map[string]any) pkg.CategoriesConfig {
	if c["default"] == nil {
		panic("must define \"default\" category")
	}
	res := pkg.CategoriesConfig{}
	set := map[string]bool{}
	for categoryName, config := range c {
		if set[categoryName] {
			panic("there are two or more appenders with the same name!")
		}
		set[categoryName] = true
		switch v := config.(type) {
		case map[string]any:
			res[categoryName] = validateCategoryConfig(v)
		default:
			panic(fmt.Sprintf("invalid type for category with name %q: required map[string]any, found %T", categoryName, v))
		}
	}
	return res
}

func validateCategoryConfig(c map[string]any) pkg.CategoryConfig {
	const (
		AppendersField = "appenders"
		LevelField     = "level"
	)
	switch v := c[AppendersField].(type) {
	case string:
		c[AppendersField] = []string{v}
	case []string:
		// correct
	case []any:
		for _, value := range v {
			switch value.(type) {
			case string:
			default:
				panic(fmt.Sprintf("appenders field type invalid. required []string, found %T", value))
			}
		}
	default:
		panic(fmt.Sprintf("appenders field type invalid. required []string, found %T", v))
	}
	switch v := c[LevelField].(type) {
	case string:
		// ok
	default:
		panic(fmt.Sprintf("level field type invalid. required string, found %T", v))
	}
	return c
}
