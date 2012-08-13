package steno

type LogLevel struct {
	name     string
	priority int
}

var LOG_OFF = NewLogLevel("off", 0)
var LOG_FATAL = NewLogLevel("fatal", 1)
var LOG_ERROR = NewLogLevel("error", 5)
var LOG_WARN = NewLogLevel("warn", 10)
var LOG_INFO = NewLogLevel("info", 15)
var LOG_DEBUG = NewLogLevel("debug", 16)
var LOG_DEBUG1 = NewLogLevel("debug1", 17)
var LOG_DEBUG2 = NewLogLevel("debug2", 18)
var LOG_ALL = NewLogLevel("all", 30)

var LEVELS = map[string]*LogLevel{
	"off":    LOG_OFF,
	"fatal":  LOG_FATAL,
	"error":  LOG_ERROR,
	"warn":   LOG_WARN,
	"info":   LOG_INFO,
	"debug":  LOG_DEBUG,
	"debug1": LOG_DEBUG1,
	"debug2": LOG_DEBUG2,
	"all":    LOG_ALL,
}

func NewLogLevel(name string, priority int) *LogLevel {
	level := new(LogLevel)

	level.name = name
	level.priority = priority

	return level
}

func lookupLevel(name string) *LogLevel {
	return LEVELS[name]
}

func (level *LogLevel) MarshalJSON() ([]byte, error) {
	return []byte(level.name), nil
}
