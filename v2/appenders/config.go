package appenders

import (
	"encoding/json"

	"github.com/alessandrofascini/log4go/v2/layouts"
)

type Config struct {
	Typo   string          `json:"type"`
	Layout *layouts.Config `json:"layout"`
	source json.RawMessage
}

func (c *Config) UnmarshalJSON(data []byte) error {
	conf := new(struct {
		Typo   string          `json:"type"`
		Layout json.RawMessage `json:"layout"`
	})
	if err := json.Unmarshal(data, conf); err != nil {
		return err
	}
	c.Typo = conf.Typo
	c.source = data
	c.Layout = new(layouts.Config)
	return c.Layout.UnmarshalJSON(conf.Layout)
}
