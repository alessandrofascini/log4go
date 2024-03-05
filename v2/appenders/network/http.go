package network

// TODO: implement me

type HttpConfig struct {
}

func (c *HttpConfig) UnmarshalJSON(data []byte) error {
	return nil
}

type Http struct {
}

func NewHttp(conf *HttpConfig) (*Http, error) {
	return nil, nil
}

func (h *Http) Write(b []byte) (int, error) {
	return 0, nil
}

func (h *Http) UnmarshalJSON(confNet []byte) error {
	return nil
}
