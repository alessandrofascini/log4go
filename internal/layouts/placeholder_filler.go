package layouts

import (
	"fmt"

	"github.com/alessandrofascini/log4go/internal/layouts/types"
)

type iPlaceholderFiller interface {
	fill(string, types.LayoutContext) any
}

type stdPlaceholderFiller struct {
	placeholders []iField
}

func (s *stdPlaceholderFiller) fill(pattern string, context types.LayoutContext) any {
	var a []any
	for _, p := range s.placeholders {
		a = append(a, p.resolve(context))
	}
	return fmt.Sprintf(pattern, a...)
}

func createStdPlaceholder(s string, occ [][]int) *stdPlaceholderFiller {
	var fields []iField
	for _, o := range occ {
		l := o[0]
		u := o[1]
		fields = append(fields, *createField(s[l:u]))
	}
	return &stdPlaceholderFiller{
		placeholders: fields,
	}
}
