package fields

import (
	"io"
	"os/user"
)

type HostnameField struct{}

func (h *HostnameField) WriteTo(w io.Writer, reader *EventReader) (n int, err error) {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	return w.Write([]byte(u.Username))
}
