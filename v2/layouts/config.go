package layouts

import (
	"encoding/json"
)

type Config struct {
	LayoutType string `json:"type"`
	Pattern    string `json:"pattern"`
	source     json.RawMessage
}

func (c *Config) UnmarshalJSON(data []byte) error {
	t := new(struct {
		Type    string `json:"type"`
		Pattern string `json:"pattern"`
	})
	if err := json.Unmarshal(data, t); err != nil {
		return err
	}
	c.LayoutType = t.Type
	c.Pattern = t.Pattern
	c.source = data
	return nil
}

// func NewConfig(src json.RawMessage) (*Config, error) {
// 	conf := &Config{}
// 	conf.source = src
// 	var err error
// 	if conf.LayoutType, err = getStringValue(src, "field"); err != nil {
// 		return nil, err
// 	}
// 	if conf.Pattern, err = getStringValue(src, "pattern"); err != nil {
// 		return nil, err
// 	}
// 	return conf, nil
// }

// func getStringValue(m sourceConfig, key string) (string, error) {
// 	if v, ok := m[key].(string); ok {
// 		return v, nil
// 	}
// 	return "", ErrMissingField
// }
