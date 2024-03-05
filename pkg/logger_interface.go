package pkg

type ILogger interface {
	AddContext(key string, value any)
	ChangeOneContext(key string, value any)
	ChangeManyContext(c map[string]any)
	RemoveContext(key string)
	ClearContext()
	Trace(args ...any)
	Debug(args ...any)
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Fatal(args ...any)
	Terminate()
}
