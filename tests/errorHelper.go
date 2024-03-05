package tests

import "os"

func WriteError(e error) {
	os.Stderr.WriteString(e.Error() + "\n")
}
