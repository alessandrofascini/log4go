package layouts

import "github.com/alessandrofascini/log4go/internal/layouts/types"

type iField interface {
	resolve(types.LayoutContext) any
}
