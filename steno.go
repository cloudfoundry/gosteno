package steno

import "encoding/json"

var config Config
var loggers map[string]Logger

func Init(c *Config) {
	config = *c

	loggers = make(map[string]Logger)

	if c.port > 0 {
		initHttpServer(c.port)
	}
}

func loggersInJson() string {
	bytes, _ := json.Marshal(loggers)
	return string(bytes)
}
