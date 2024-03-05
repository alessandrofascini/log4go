package layouts

import (
	"fmt"
	"os/user"
	"regexp"
	"strconv"
	"strings"

	"github.com/alessandrofascini/log4go/pkg"
)

// NewLayout This function create a layout starting from a LayoutConfig. Returns a Layout
func NewLayout(config pkg.LayoutConfig) pkg.Layout {
	t := config.GetType()
	switch t {
	case "basic":
		config["pattern"] = "[%d] [%p] %c - %m"
		return patternLayout(&config)
	case "messagePassThrough":
		config["pattern"] = "%m"
		return patternLayout(&config)
	case "", "colored", "coloured":
		config["pattern"] = "%[[%d] [%p] %c%] - %m"
		return patternLayout(&config)
	case "pattern":
		return patternLayout(&config)
	}
	if externalLayoutPool[t] != nil {
		return (*externalLayoutPool[t])(&config)
	}
	panic(fmt.Sprintf("unsopported layout type: %q", t))
}

func patternLayout(config *pkg.LayoutConfig) pkg.Layout {
	// Get Pattern
	var p = (*config)["pattern"]
	if p == nil {
		panic("cannot resolve layout pattern without a pattern!")
	}
	var pattern = ""
	switch p := p.(type) {
	default:
		panic("pattern must be a string")
	case string:
		pattern = p
	}

	format, functions := preparePattern(pattern, config)

	return func(event pkg.LoggingEvent) string {
		var args []any
		for _, f := range functions {
			args = append(args, f(event))
		}
		return fmt.Sprintf(format, args...)
	}
}

var patternRegex = regexp.MustCompile("%((-\\d+|\\d*)?(\\.-\\d+|\\.\\d*)?)?.(\\{[^}]*})?")

func preparePattern(pattern string, config *pkg.LayoutConfig) (string, []func(event pkg.LoggingEvent) any) {
	occ := patternRegex.FindAllStringIndex(pattern, -1)
	prev := 0
	strBuilder := strings.Builder{}
	var filler []func(event pkg.LoggingEvent) any
	for _, v := range occ {
		start := v[0]
		end := v[1]
		strBuilder.WriteString(pattern[prev:start])
		strBuilder.WriteString("%v")
		filler = append(filler, getFillerFunction(pattern[start:end], config))
		prev = end
	}
	strBuilder.WriteString(pattern[prev:])
	// for avoiding more logs on the same line
	strBuilder.WriteString("\n")
	return strBuilder.String(), filler
}

const (
	paddingIndex    = iota + 1
	truncationIndex = iota + 2
	fieldIndex
	formatIndex = iota + 3
)

var fieldRgx = regexp.MustCompile("%(-\\d+|\\d*)?(\\.(-\\d+|\\d*))?(.)(\\{(.*)})?$")
var fieldFunctions = []func(string, *pkg.LayoutConfig) func(pkg.LoggingEvent) any{
	// %r - time in toLocaleTimeString format
	114: func(string, *pkg.LayoutConfig) func(pkg.LoggingEvent) any {
		return func(pkg.LoggingEvent) any {
			return pkg.GetDateByFormat(pkg.LocaleTime)
		}
	},
	// %p log level
	112: func(string, *pkg.LayoutConfig) func(pkg.LoggingEvent) any {
		return func(event pkg.LoggingEvent) any {
			return event.Level.LevelStr
		}
	},
	// %c log category ~ category name
	99: func(string, *pkg.LayoutConfig) func(pkg.LoggingEvent) any {
		return func(event pkg.LoggingEvent) any {
			return event.CategoryName
		}
	},
	// %h hostname
	104: func(string, *pkg.LayoutConfig) func(pkg.LoggingEvent) any {
		return func(pkg.LoggingEvent) any {
			u, err := user.Current()
			if err != nil {
				panic(err)
			}
			return u.Username
		}
	},
	// %m log data
	109: func(string, *pkg.LayoutConfig) func(pkg.LoggingEvent) any {
		return func(event pkg.LoggingEvent) any {
			format := strings.Repeat("%v ", len(event.Data))
			format = format[0 : len(format)-1]
			return fmt.Sprintf(format, event.Data...)
		}
	},
	// %d date formatted
	100: func(format string, _ *pkg.LayoutConfig) func(pkg.LoggingEvent) any {
		return func(pkg.LoggingEvent) any {
			switch format {
			case "", "ISO8601": // default case
				return pkg.GetDateByFormat(pkg.ISO8601)
			case "ISO8601_WITH_TZ_OFFSET":
				return pkg.GetDateByFormat(pkg.Iso8601WithTzOffset)
			case "ABSOLUTETIME":
				return pkg.GetDateByFormat(pkg.AbsoluteTime)
			case "DATETIME":
				return pkg.GetDateByFormat(pkg.Datetime)
			}
			// not found standard format
			return pkg.GetDateByFormat(format)
		}
	},
	// %% % - for when you want a literal % in your output
	37: func(string, *pkg.LayoutConfig) func(pkg.LoggingEvent) any {
		return func(pkg.LoggingEvent) any {
			return "%"
		}
	},
	// %n newline
	110: func(string, *pkg.LayoutConfig) func(pkg.LoggingEvent) any {
		return func(pkg.LoggingEvent) any {
			return "\n"
		}
	},
	// %z process id
	122: func(string, *pkg.LayoutConfig) func(pkg.LoggingEvent) any {
		return func(event pkg.LoggingEvent) any {
			return event.Pid
		}
	},
	// %x{<tokenname>} add dynamic tokens to your log. Tokens are specified in the tokens parameter.
	120: func(token string, config *pkg.LayoutConfig) func(pkg.LoggingEvent) any {
		var tokens interface{}
		tokens = (*config)["tokens"]
		if tokens == nil {
			tokens = make(pkg.LayoutTokens)
		}
		return func(pkg.LoggingEvent) any {
			switch v := tokens.(type) {
			case map[string]any:
				return v[token]
			}
			return "nil"
		}
	},
	// %X{<tokenname>} add values from the Log4GoLogger context. Tokens are keys into the context values.
	88: func(token string, _ *pkg.LayoutConfig) func(pkg.LoggingEvent) any {
		return func(event pkg.LoggingEvent) any {
			if event.Context[token] == nil {
				return "nil"
			}
			var v interface{}
			v = event.Context[token]
			switch v := v.(type) {
			case func(pkg.LoggingEvent) any:
				return v(event)
			case func() any:
				return v()
			}
			return v
		}
	},
	// %[ start a coloured block (colour will be taken from the log level, similar to colouredLayout)
	91: func(string, *pkg.LayoutConfig) func(pkg.LoggingEvent) any {
		return func(event pkg.LoggingEvent) any {
			return event.Level.Color.Code
		}
	},
	// %] end a coloured block
	93: func(string, *pkg.LayoutConfig) func(pkg.LoggingEvent) any {
		return func(pkg.LoggingEvent) any {
			return pkg.ColorReset.Code
		}
	},
}

func getFillerFunction(s string, config *pkg.LayoutConfig) func(event pkg.LoggingEvent) any {
	match := fieldRgx.FindAllStringSubmatch(s, -1)[0]
	padding, _ := strconv.Atoi(match[paddingIndex])       // on error padding => 0
	truncation, _ := strconv.Atoi(match[truncationIndex]) // on error truncation => 0
	fieldName := match[fieldIndex][0]                     // a b c p d
	format := match[formatIndex]                          // { }

	f := fieldFunctions[fieldName](format, config)

	// truncation and padding
	if truncation != 0 && padding != 0 {
		return func(event pkg.LoggingEvent) any {
			t := fmt.Sprintf("%v", f(event))
			t = truncate(t, truncation)
			t = paddingString(t, padding)
			return t
		}
	}

	// truncation
	if truncation != 0 {
		return func(event pkg.LoggingEvent) any {
			t := fmt.Sprintf("%v", f(event))
			return truncate(t, truncation)
		}
	}

	// padding
	if padding != 0 {
		return func(event pkg.LoggingEvent) any {
			t := fmt.Sprintf("%v", f(event))
			return paddingString(t, padding)
		}
	}
	return f
}
