package steno

var config Config
var loggers map[string]*Logger

func Init(c *Config) {
	config = *c

	loggers = make(map[string]*Logger)
}
