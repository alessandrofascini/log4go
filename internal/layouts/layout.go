package layouts

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/alessandrofascini/log4go/internal/layouts/types"
)

type layout struct {
	lType   types.LayoutType
	pattern string
	filler  iPlaceholderFiller
}

func (l *layout) Fill(ctx types.LayoutContext) any {
	return l.filler.fill(l.pattern, ctx)
}

func createLayout(layoutType types.LayoutType, pattern types.PatternConfig) *iLayout {
	strPattern := string(pattern)
	fieldRegex := regexp.MustCompile("%((-\\d+|\\d*)?(\\.-\\d+|\\.\\d*)?)?.(\\{[^}]*})?")

	placeholderOccurrences := fieldRegex.FindAllStringIndex(strPattern, -1)

	var l iLayout = &layout{
		lType:   layoutType,
		pattern: createPattern(strPattern, placeholderOccurrences),
		filler:  createStdPlaceholder(strPattern, placeholderOccurrences),
	}
	fmt.Println(l)
	return &l
}

func createPattern(origin string, occ [][]int) string {
	builder := &strings.Builder{}
	prev := 0
	for _, o := range occ {
		l := o[0]
		builder.WriteString(origin[prev:l])
		builder.WriteString("%v")
		prev = o[1]
	}
	builder.WriteString(origin[prev:])
	return builder.String()
}
