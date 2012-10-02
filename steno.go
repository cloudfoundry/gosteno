package steno

import (
	"encoding/json"
	"sync"
)

// Global configs
var config Config

// loggersMutex protects accesses to loggers and regexp
var loggersMutex sync.Mutex

// loggers only saves BaseLogger
var loggers = make(map[string]*BaseLogger)

func Init(c *Config) {
	config = *c

	if config.Level == nil {
		config.Level = LOG_INFO
	}
	if config.Codec == nil {
		config.Codec = JSON_CODEC
	}
	if config.Sinks == nil {
		config.Sinks = []Sink{}
	}

	for _, sink := range config.Sinks {
		if sink.GetCodec() == nil {
			sink.SetCodec(config.Codec)
		}
	}

	if config.Port > 0 {
		initHttpServer(config.Port)
	}

	for name, _ := range loggers {
		loggers[name] = nil
	}
}

func NewLogger(name string) Logger {
	loggersMutex.Lock()
	defer loggersMutex.Unlock()

	logger := loggers[name]

	if logger == nil {
		baseLogger := new(BaseLogger)

		baseLogger.name = name
		baseLogger.sinks = config.Sinks
		baseLogger.level = computeLevel(name)

		logger = baseLogger
		loggers[name] = logger
	}

	return logger
}

func loggersInJson() string {
	bytes, _ := json.Marshal(loggers)
	return string(bytes)
}
