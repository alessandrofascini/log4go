package layouts

import (
	"fmt"

	"github.com/alessandrofascini/log4go/internal/layouts/types"
	"github.com/alessandrofascini/log4go/internal/layouts/utils"
)

type paddingDecoratorField struct {
	wrapper iField
	pad     int
}

func (l *paddingDecoratorField) resolve(cxt types.LayoutContext) any {
	prev := l.wrapper.resolve(cxt)
	s := fmt.Sprintf("%v", prev)
	return utils.PaddingString(s, l.pad)
}
