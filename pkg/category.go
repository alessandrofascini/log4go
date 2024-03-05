package pkg

import (
	"fmt"
	"strings"
)

type CategoryConfig map[string]any

func (c CategoryConfig) GetAssociatedAppenders() []string {
	v := c["appenders"]
	switch a := v.(type) {
	case string:
		return []string{a}
	case []string:
		return a
	case []any:
		var res []string
		for _, value := range a {
			switch v := value.(type) {
			case string:
				res = append(res, v)
			default:
				panic(fmt.Sprintf("appenders field type invalid. required []string, found %T", value))
			}
		}
		return res
	}
	panic(fmt.Sprintf("unsupported appenders pkg %T", v))
}

func (c CategoryConfig) GetLevelValue() uint {
	v, ok := c["level"]
	if !ok {
		panic(fmt.Errorf("cannot find field %q in category config %v", "level", c))
	}
	s := fmt.Sprintf("%v", v)
	return GetLevelValueByName(strings.ToUpper(s))
}

func (c CategoryConfig) GetName() string {
	v, ok := c["name"]
	if !ok {
		panic(fmt.Errorf("cannot find field %q in category config %v", "name", c))
	}
	return fmt.Sprintf("%v", v)
}

func (c CategoryConfig) SetName(n string) {
	c["name"] = n
}

type CategoriesConfig map[string]CategoryConfig
