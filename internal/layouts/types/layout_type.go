package types

type LayoutType uint8

const (
	LayoutTypeDefault LayoutType = iota
	LayoutTypeBasic
	LayoutTypeColoured
	LayoutTypeMessagePassThrough
	LayoutTypePattern
)

var layoutTypesName = []string{"default", "basic", "coloured", "messagePassThrough", "pattern"}

func (lt LayoutType) String() string {
	if lt > LayoutTypePattern {
		panic("unknown layout type")
	}
	return layoutTypesName[lt]
}

func (lt LayoutType) IsValid() bool {
	return lt <= LayoutTypePattern
}
