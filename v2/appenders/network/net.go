package network

import (
	"encoding/json"
	"net"
)

type NetConfig struct {
	Port     string `json:"port"`
	Host     string `json:"host"`
	Protocol string `json:"protocol"`
}

func (n *NetConfig) UnmarshalJSON(data []byte) error {
	t := new(struct {
		Port     string `json:"port"`
		Host     string `json:"host"`
		Protocol string `json:"protocol"`
	})
	if err := json.Unmarshal(data, t); err != nil {
		return err
	}
	n.Protocol = t.Protocol
	n.Host = t.Host
	n.Port = t.Port
	return nil
}

// Net implements io.WriteCloser and json.UnmarshalJSON
type Net struct {
	conn              net.Conn
	protocol, address string
}

func NewNet(conf *NetConfig) (*Net, error) {
	return &Net{nil, conf.Protocol, net.JoinHostPort(conf.Host, conf.Port)}, nil
}

func (n *Net) Write(b []byte) (int, error) {
	if n.conn == nil {
		var err error
		if n.conn, err = net.Dial(n.protocol, n.address); err != nil {
			return 0, err
		}
	}
	return n.conn.Write(b)
}

func (n *Net) Close() error {
	return n.conn.Close()
}

func (n *Net) UnmarshalJSON(confNet []byte) error {
	conf := &NetConfig{}
	if err := conf.UnmarshalJSON(confNet); err != nil {
		return err
	}
	var err error
	n, err = NewNet(conf)
	return err
}
