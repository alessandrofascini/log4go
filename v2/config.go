package v2

import (
	"encoding/json"

	"github.com/alessandrofascini/log4go/v2/appenders"
	"github.com/alessandrofascini/log4go/v2/categories"
	"github.com/alessandrofascini/log4go/v2/internal/config"
)

type Config struct {
	Appenders     appenders.Config  `json:"appenders"`
	Categories    categories.Config `json:"categories"`
	Configuration config.Config     `json:"configuration"`
}

func (c *Config) UnmarshalJSON(data []byte) error {
	type configJSON struct {
		Appenders     json.RawMessage `json:"appenders"`
		Categories    json.RawMessage `json:"categories"`
		Configuration json.RawMessage `json:"configuration"`
	}
	conf := &configJSON{}
	if err := json.Unmarshal(data, conf); err != nil {
		return err
	}
	// unmarshal appenders
	if err := c.Appenders.UnmarshalJSON(conf.Appenders); err != nil {
		return err
	}
	// unmarshal categories

	// unmarshal configuration

	return nil
}
