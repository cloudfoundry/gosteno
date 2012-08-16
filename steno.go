package steno

import (
	"encoding/json"
	"regexp"
	"sync"
)

// Global configs
var config Config

// loggersMutex protects accesses to loggers and regexp
var loggersMutex sync.Mutex

// loggers only saves BaseLogger
var loggers = make(map[string]*BaseLogger)

// loggerRegexp* used to match log name and log level
var loggerRegexp *regexp.Regexp
var loggerRegexpLevel *LogLevel

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

func SetLoggerRegexp(pattern string, level *LogLevel) error {
	loggersMutex.Lock()
	defer loggersMutex.Unlock()

	clearLoggerRegexp()
	return setLoggerRegexp(pattern, level)
}

func ClearLoggerRegexp() {
	loggersMutex.Lock()
	defer loggersMutex.Unlock()

	clearLoggerRegexp()
}

func setLoggerRegexp(pattern string, level *LogLevel) error {
	regExp, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	// If here, Logger regexp is valid
	loggerRegexp = regExp
	loggerRegexpLevel = level
	for name, logger := range loggers {
		if loggerRegexp.MatchString(name) {
			logger.level = loggerRegexpLevel
		}
	}

	return nil
}

func clearLoggerRegexp() {
	if loggerRegexp == nil {
		return
	}

	for name, logger := range loggers {
		if loggerRegexp.MatchString(name) {
			logger.level = config.Level
		}
	}

	loggerRegexp = nil
	loggerRegexpLevel = nil
}

func computeLevel(name string) *LogLevel {
	if loggerRegexp != nil && loggerRegexp.MatchString(name) {
		return loggerRegexpLevel
	}

	return config.Level
}

func loggersInJson() string {
	bytes, _ := json.Marshal(loggers)
	return string(bytes)
}
