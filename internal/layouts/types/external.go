package types

type TokenFn func() any
type Tokens map[string]TokenFn

type LayoutConfiguration interface {
	GetType() LayoutType
	GetPattern() PatternConfig
}

type LayoutContext interface {
	GetLevel() string
	GetCategoryName() string
	GetLogData() string
	GetTokens() Tokens
}
