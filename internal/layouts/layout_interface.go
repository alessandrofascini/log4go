package layouts

import "github.com/alessandrofascini/log4go/internal/layouts/types"

type iLayout interface {
	Fill(ctx types.LayoutContext) any
}
