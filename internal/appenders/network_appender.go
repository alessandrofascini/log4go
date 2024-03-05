package appenders

import (
	"log"
	"net"

	"github.com/alessandrofascini/log4go/pkg"
)

func NetworkAppender(config pkg.AppenderConfig, layout *pkg.Layout) pkg.Appender {
	var protocol, host, port string
	var ok bool
	if protocol, ok = config["protocol"].(string); !ok {
		panic("missing protocol")
	}
	if host, ok = config["host"].(string); !ok {
		panic("missing host")
	}
	if port, ok = config["port"].(string); !ok {
		panic("missing port")
	}
	return func(event pkg.LoggingEvent) {
		conn, err := net.Dial(protocol, host+":"+port)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		if _, err := conn.Write([]byte((*layout)(event))); err != nil {
			log.Fatal(err)
		}
	}
}
