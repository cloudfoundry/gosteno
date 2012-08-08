package steno

type Logger struct {
	name string
	sinks []Sink
	level *LogLevel
}

func NewLogger(name string) *Logger {
	logger := loggers[name]

	if logger == nil {
		logger = new(Logger)

		logger.name = name
		logger.sinks = config.sinks

		level := lookupLevel(config.level)
		logger.level = level

		loggers[name] = logger
	}


	return logger
}

func (l *Logger) Log(level string, message string) {
	if !l.active(level) {
		return
	}

	for _, sink := range l.sinks {
		record := NewRecord(lookupLevel(level), message)

		sink.AddRecord(record)
		sink.Flush()
	}
}

// TODO: more functions needed
func (l *Logger) Info(message string) {
	l.Log("info", message)
}

func (l *Logger) Debug(message string) {
	l.Log("debug", message)
}

func (l *Logger) active(level string) bool {
	theLevel := lookupLevel(level)
	if theLevel == nil {
		return false
	}

	if l.level.priority >= theLevel.priority {
		return true
	}

	return false
}

// For testing
func NumLogger() int {
	return len(loggers)
}
