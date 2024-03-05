package errors

import (
	"fmt"
	"os"
)

func WriteError(e error) {
	os.Stderr.WriteString(e.Error() + "\n")
}

func WriteErrorf(format string, a ...any) {
	os.Stderr.WriteString(fmt.Sprintf(format+"\n", a))
}
