package steno

import (
	"regexp"
)

// loggerRegexp* used to match log name and log level
var loggerRegexp *regexp.Regexp
var loggerRegexpLevel *LogLevel

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
