package utils

import "runtime"

var ColorReset = "\033[0m"
var ColorRed = "\033[31m"
var ColorGreen = "\033[32m"
var ColorYellow = "\033[33m"
var ColorBlue = "\033[34m"
var ColorMagenta = "\033[35m"
var ColorCyan = "\033[36m"
var ColorGray = "\033[37m"
var ColorWhite = "\033[97m"

func init() {
	if runtime.GOOS == "windows" {
		ColorReset = ""
		ColorRed = ""
		ColorGreen = ""
		ColorYellow = ""
		ColorBlue = ""
		ColorMagenta = ""
		ColorCyan = ""
		ColorGray = ""
		ColorWhite = ""
	}
}
