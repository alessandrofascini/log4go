package appenders

import (
	"fmt"
	"testing"
)

func TestConfigUnmarshalJSON(t *testing.T) {
	jsonBlob := []byte(`{
			"type": "net",
			"layout": {
				"type": "basic",
				"pattern": "%m"
			},
			"protocol": "tcp",
			"port": "8080",
			"host": "localhost"
		}`)
	conf := new(Config)
	conf.UnmarshalJSON(jsonBlob)

	fmt.Println("type:", conf.Typo)
	fmt.Println("layout:", conf.Layout)
}
