package tests

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestUnmarshelingJSON(t *testing.T) {
	var jsonBlob = []byte(`{"appenders":{"a":5}}`)
	type Config struct {
		Appenders  json.RawMessage `json:"appenders"`
		Categories json.RawMessage `json:"categories"`
	}
	var conf Config
	err := json.Unmarshal(jsonBlob, &conf)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v\n", conf)
}
