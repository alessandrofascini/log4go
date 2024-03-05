package entities

type LogLevel int

const (
	ALL LogLevel = iota
	TRACE
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	OFF
)

var levelNames = [8]string{"ALL", "TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL", "OFF"}

func (l LogLevel) String() string {
	if l < ALL || l > OFF {
		panic("Unknown level")
	}
	return levelNames[l]
}

func (l LogLevel) IsValid() bool {
	switch l {
	case ALL, TRACE, DEBUG, INFO, WARN, ERROR, FATAL, OFF:
		return true
	}
	return false
}
