package appenders

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/alessandrofascini/log4go/pkg"

	"github.com/alessandrofascini/log4go/internal/errors"
)

/*
	Http Request config
	method: string
	url: string
	proto: string
	header http.Header
*/

func HttpAppender(config pkg.AppenderConfig, layout *pkg.Layout) pkg.Appender {
	client := &http.Client{}
	var method, url, proto string
	var ok bool
	if method, ok = config["method"].(string); !ok {
		method = "GET"
	}
	if url, ok = config["url"].(string); !ok {
		panic("missing url")
	}
	if proto, ok = config["proto"].(string); !ok {
		proto = ""
	}
	header := parseHeader(config)
	return func(event pkg.LoggingEvent) {
		v := (*layout)(event)
		body := strings.NewReader(v)
		req, err := http.NewRequest(method, url, body)
		if proto != "" {
			req.Proto = proto
		}
		req.Header = header
		req.ContentLength = int64(len(v))
		if err != nil {
			errors.WriteError(err)
			return
		}
		if _, err := client.Do(req); err != nil {
			errors.WriteError(err)
		}
	}
}

func parseHeader(config pkg.AppenderConfig) http.Header {
	header, ok := config["header"]
	if !ok {
		return http.Header{}
	}
	switch h := header.(type) {
	case http.Header:
		return h
	case map[string][]string:
		return http.Header(h)
	case map[string]string:
		return parseMapStringStringToHeader(h)
	case map[string]any:
		return parseMapStringAnyToHeader(h)
	}
	return http.Header{}
}

func parseMapStringAnyToHeader(m map[string]any) http.Header {
	header := http.Header{}
	for key, value := range m {
		switch v := value.(type) {
		case string:
			header.Add(key, v)
		case []string:
			for i := range v {
				header.Add(key, v[i])
			}
		case []any:
			for i := range v {
				header.Add(key, fmt.Sprintf("%v", v[i]))
			}
		}
	}
	return header
}

func parseMapStringStringToHeader(m map[string]string) http.Header {
	header := http.Header{}
	for key, value := range m {
		header.Add(key, value)
	}
	return header
}
