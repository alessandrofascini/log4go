package utils

import "fmt"

func PaddingString(s string, pad int) string {
	format := fmt.Sprintf("%%%ds", pad)
	return fmt.Sprintf(format, s)
}
