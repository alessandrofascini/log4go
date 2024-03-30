package v2_test

type LevelWriter interface {
	Trace()
	Debug()
	Info()
	Error()
	Warning()
	Fatal()
}
