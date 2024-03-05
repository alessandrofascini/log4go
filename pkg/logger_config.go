package pkg

type LoggerConfig map[string]any

func (lc LoggerConfig) GetAppenders() AppendersConfig {
	switch v := lc["appenders"].(type) {
	case AppendersConfig:
		return v
	}
	panic("LoggerConfig.getAppenders() must return AppendersConfig")
}

func (lc LoggerConfig) GetCategories() CategoriesConfig {
	switch v := lc["categories"].(type) {
	case CategoriesConfig:
		return v
	}
	panic("LoggerConfig.getCategories() must return CategoriesConfig")
}
