package layouts

import (
	"fmt"

	"github.com/alessandrofascini/log4go/internal/layouts/types"
	"github.com/alessandrofascini/log4go/internal/layouts/utils"
)

type truncateDecoratorField struct {
	wrapper iField
	trunc   int
}

func (l *truncateDecoratorField) resolve(cxt types.LayoutContext) any {
	prev := l.wrapper.resolve(cxt)
	s := fmt.Sprintf("%v", prev)
	return utils.Truncate(s, l.trunc)
}
