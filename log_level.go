package steno

type LogLevel struct {
	name string
	priority int
}

var LEVELS = map[string]*LogLevel{
	"off"    : NewLogLevel("off",     0),
	"fatal"  : NewLogLevel("fatal",   1),
	"error"  : NewLogLevel("error",   5),
	"warn"   : NewLogLevel("warn",   10),
	"info"   : NewLogLevel("info",   15),
	"debug"  : NewLogLevel("debug",  16),
	"debug1" : NewLogLevel("debug1", 17),
	"debug2" : NewLogLevel("debug2", 18),
	"all"    : NewLogLevel("all",    30),
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
