# Log4Go

---

**Table of Contents**

- [Terminology](#terminology)
- [API](#API)
- [Configuration](#configuration)
  - [Appenders](#appenders)
  - [Categories](#categories)
  - [Layouts](#layouts)
- [Next Steps](#next-steps)

---

## Terminology

---

**Level** - a log level is the severity or priority of a log event (debug, info, etc).
Whether an _appender_ will see the event or not is determined by the _category_’s level.
If this is less than or equal to the event’s level, it will be sent to the category’s appender(s).

**Category** - a label for grouping log events.
This can be based on module (e.g. ‘auth’, ‘payment’, ‘http’), or anything you like.
Log events with the same category will go to the same appenders.
The category for log events is defined
when you get a Logger from Log4Go (**Log4Go.GetLogger('some-category')**).

**Appender** - appenders are responsible for output of log events.
They may write events to files, send emails, store them in a database, or anything.
Most appenders use _layouts_ to serialise the events to strings for output.

**Logger** - this is your code’s main interface with Log4Go.
A logger instance may have an optional _category_, defined when you create the instance.
Loggers provide the **info, debug, error**, ... functions that create _LogEvents_ and pass them on to appenders.

**Layout** - a function for converting a _LogEvent_ into a string representation.
Log4Go comes with a few different implementations: **basic**, **coloured**, and a more configurable pattern based layout.

**LogEvent** - a log event has a timestamp, a level, and optional category, data, and context properties.
When you call logger.Info('cheese value:', edam) the logger will create a log event with the timestamp of now,
a _level_ of INFO, a category that was chosen when the logger was created,
and a data array with two values (the string ‘cheese value:’, and the object ‘edam’),
along with any context data that was added to the logger.

# API

## Configuration

---

#### configuration - Log4Go.Configure(map\[string\]any | string)

This method allows you to configure the logger.
You can pass a JSON (such as map[string]any), or, if you prefer, you can pass a string as the path to a json file

#### Configuration Object

Properties:

- appenders (map\[string\]any) - a map of named appenders (string) to appender definitions (object); appender definitions must have a property type (string) - other properties depend on the appender type.
- categories (map\[string\]any) - a map of named categories (string) to category definitions (object).
  You must define the default category which is used for all log events that do not match a specific category.
  Category definitions have two properties:
  - appenders (array of strings) - the list of appender names to be used for this category. A category must have at least one appender.
  - level (string, case insensitive) - the minimum log level that this category will send to the appenders. For example, if set to ‘error’ then the appenders will only receive log events of level ‘error’, ‘fatal’, ‘mark’ - log events of ‘info’, ‘warn’, ‘debug’, or ‘trace’ will be ignored.

#### configured - Log4Go.IsConfigured() bool

IsConfigured method call returns a boolean on whether Log4Go.Configure() was successfully called previously.
Implicit Log4Go.Configure() call by Log4Go.GetLogger() is will also affect this value.

#### Loggers - Log4Go.GetLogger(_category-name_ string) Logger

To support the minimalist usage, this function will implicitly call Log4Go.Configure() with the default configurations if it hasn’t been configured before.

#### @Type Logger

```go
package pkg

type ILog4GoLogger interface {
	AddContext(key string, value any)
	ChangeOneContext(key string, value any)
	ChangeManyContext(c map[string]any)
	RemoveContext(key string)
	ClearContext()
	Trace(args ...any)
	Debug(args ...any)
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Fatal(args ...any)
	Terminate()
}

```

## Appenders

---

### Type of Appenders

- dateFile
- file
- fileSync
- stderr
- stdout

## Categories

---

## Layouts

---

### Type of Layouts

- basic
- coloured
- messagePassThrough
- pattern
