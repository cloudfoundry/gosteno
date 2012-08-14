package steno

import "encoding/json"

var config Config
var loggers = make(map[string]Logger)

func Init(c *Config) {
	config = *c

	if config.Level == nil {
		config.Level = LOG_INFO
	}
	if config.Codec == nil {
		config.Codec = JSON_CODEC
	}
	if config.Sinks == nil || len(config.Sinks) == 0 {
		panic("Cannot init with no sinks")
	}

	for _, sink := range config.Sinks {
		sink.SetCodec(config.Codec)
	}

	if config.Port > 0 {
		initHttpServer(config.Port)
	}
}

func loggersInJson() string {
	bytes, _ := json.Marshal(loggers)
	return string(bytes)
}
