package pkg

type Color struct {
	Name string
	Code string
}

var (
	ColorReset   = &Color{"reset", "\033[0m"}
	ColorRed     = &Color{"red", "\033[31m"}
	ColorGreen   = &Color{"green", "\033[32m"}
	ColorYellow  = &Color{"yellow", "\033[33m"}
	ColorBlue    = &Color{"Blue", "\033[34m"}
	ColorMagenta = &Color{"magenta", "\033[35m"}
	ColorCyan    = &Color{"cyan", "\033[36m"}
	//ColorGray    = &Color{"gray", "\033[37m"}
	//ColorWhite   = &Color{"white", "\033[97m"}
)
