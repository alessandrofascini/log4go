package entities

type LogEntity struct {
	level        LogLevel
	categoryName string
	message      string
}

func (e *LogEntity) SetLevel(level LogLevel) {
	e.level = level
}

func (e *LogEntity) SetCategoryName(categoryName string) {
	e.categoryName = categoryName
}

func (e *LogEntity) SetMessage(message string) {
	e.message = message
}

func (e *LogEntity) GetLevel() LogLevel {
	return e.level
}

func (e *LogEntity) GetCategoryName() string {
	return e.categoryName
}

func (e *LogEntity) GetMessage() string {
	return e.message
}
