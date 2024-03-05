package types

type PatternConfig string

const (
	PatternConfDefault            = "[%d] [%p] %c - %m"
	PatternConfMessagePassThrough = "%m"
	PatternConfColoured           = "%[[%d] [%p] %c%] - %m"
)
