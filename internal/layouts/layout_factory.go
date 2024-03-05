package layouts

import "github.com/alessandrofascini/log4go/internal/layouts/types"

func LayoutFactory(config types.LayoutConfiguration) iLayout {
	t := types.LayoutTypeDefault
	var p types.PatternConfig
	switch config.GetType() {
	case types.LayoutTypeDefault, types.LayoutTypeBasic:
		t = types.LayoutTypeBasic
		p = types.PatternConfDefault
		break
	case types.LayoutTypeMessagePassThrough:
		t = types.LayoutTypeMessagePassThrough
		p = types.PatternConfMessagePassThrough
		break
	case types.LayoutTypePattern:
		t = types.LayoutTypePattern
		p = config.GetPattern()
		break
	case types.LayoutTypeColoured:
		t = types.LayoutTypePattern
		p = types.PatternConfColoured
	default:
		panic("unknown config type")
	}
	return *createLayout(t, p)
}
