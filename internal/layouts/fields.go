package layouts

import (
	"os"
	"os/user"
	"regexp"
	"strconv"

	"github.com/alessandrofascini/log4go/internal/layouts/types"

	"github.com/alessandrofascini/log4go/internal/layouts/utils"
)

type fieldLocaleTime struct {
}

func (f *fieldLocaleTime) resolve(types.LayoutContext) any {
	return utils.GetStringDate(utils.LocaleTime)
}

type fieldLogLevel struct {
}

func (f *fieldLogLevel) resolve(context types.LayoutContext) any {
	return context.GetLevel()
}

type fieldLogCategoryName struct {
}

func (f *fieldLogCategoryName) resolve(context types.LayoutContext) any {
	return context.GetCategoryName()
}

type fieldHostname struct {
}

func (f *fieldHostname) resolve(types.LayoutContext) any {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	return u.Username
}

type fieldLogData struct {
}

func (f *fieldLogData) resolve(context types.LayoutContext) any {
	return context.GetLogData()
}

type fieldPercent struct {
}

func (f *fieldPercent) resolve(types.LayoutContext) any {
	return "%"
}

type fieldNewLine struct {
}

func (f *fieldNewLine) resolve(types.LayoutContext) any {
	return "\n"
}

type fieldProcessId struct{}

func (f *fieldProcessId) resolve(types.LayoutContext) any {
	return os.Getppid()
}

type fieldColorTextOpen struct {
}

func (f *fieldColorTextOpen) resolve(cxt types.LayoutContext) any {
	switch cxt.GetLevel() {
	case "TRACE":
		return utils.ColorBlue
	case "DEBUG":
		return utils.ColorCyan
	case "INFO":
		return utils.ColorGreen
	case "WARN":
		return utils.ColorYellow
	case "ERROR":
		return utils.ColorRed
	case "FATAL":
		return utils.ColorMagenta
	case "OFF", "ALL":
		return utils.ColorGray
	}
	return utils.ColorWhite
}

type fieldColorTextClose struct {
}

func (f *fieldColorTextClose) resolve(types.LayoutContext) any {
	return utils.ColorReset
}

type fieldToken struct {
	tokenKey string
}

func (f *fieldToken) resolve(cxt types.LayoutContext) any {
	tokens := cxt.GetTokens()
	function := tokens[f.tokenKey]
	if function != nil {
		return function()
	}
	return ""
}

func createField(s string) *iField {
	rgx := regexp.MustCompile("%(-\\d+|\\d*)?(\\.(-\\d+|\\d*))?(.)(\\{(.*)})?$")
	match := rgx.FindAllStringSubmatch(s, -1)[0]
	padding, _ := strconv.Atoi(match[paddingIndex])
	truncation, _ := strconv.Atoi(match[truncationIndex])
	fieldName := match[fieldIndex][0]
	format := match[formatIndex]
	var newField iField
	switch fieldName {
	case 'r':
		newField = new(fieldLocaleTime)
		break
	case 'p':
		newField = new(fieldLogLevel)
		break
	case 'c':
		newField = new(fieldLogCategoryName)
		break
	case 'h':
		newField = new(fieldHostname)
		break
	case 'm':
		newField = new(fieldLogData)
		break
	case 'd':
		newField = createLeafDate(format)
		break
	case '%':
		newField = new(fieldPercent)
		break
	case 'n':
		newField = new(fieldNewLine)
		break
	case 'z':
		newField = new(fieldProcessId)
		break
	case '[':
		newField = new(fieldColorTextOpen)
		break
	case ']':
		newField = new(fieldColorTextClose)
		break
	case 'x':
		newField = &fieldToken{tokenKey: format}
		break
	}
	if truncation != 0 {
		newField = &truncateDecoratorField{
			wrapper: newField,
			trunc:   truncation,
		}
	}
	if padding != 0 {
		newField = &paddingDecoratorField{
			wrapper: newField,
			pad:     padding,
		}
	}
	return &newField
}
