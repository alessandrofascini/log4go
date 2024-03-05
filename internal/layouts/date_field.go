package layouts

import (
	"github.com/alessandrofascini/log4go/internal/layouts/types"

	"github.com/alessandrofascini/log4go/internal/layouts/utils"
)

const (
	FieldDateFormatDefault             = ""
	FieldDateFormatISO8601             = "ISO8601"
	FieldDateFormatISO8601WithTzOffset = "ISO8601_WITH_TZ_OFFSET"
	FieldDateFormatAbsoluteTime        = "ABSOLUTETIME"
	FieldDateFormatDateTime            = "DATETIME"
)

type fieldDate struct {
	format string
}

func (f *fieldDate) resolve(types.LayoutContext) any {
	return utils.GetStringDate(f.format)
}

func createLeafDate(format string) *fieldDate {
	switch format {
	case FieldDateFormatDefault, FieldDateFormatISO8601:
		return &fieldDate{format: utils.ISO8601}
	case FieldDateFormatISO8601WithTzOffset:
		return &fieldDate{format: utils.Iso8601WithTzOffset}
	case FieldDateFormatAbsoluteTime:
		return &fieldDate{format: utils.AbsoluteTime}
	case FieldDateFormatDateTime:
		return &fieldDate{format: utils.Datetime}
	}
	return &fieldDate{format: format}
}
