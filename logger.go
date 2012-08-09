package steno

type Logger interface {
	Log(level *LogLevel, message string, data map[string]string)
	Fatal(message string)
	Error(message string)
	Warn(message string)
	Info(message string)
	Debug(message string)
	Debug1(message string)
	Debug2(message string)
}

type BaseLogger struct {
	name  string
	sinks []Sink
	level *LogLevel
}

func NewLogger(name string) Logger {
	logger := loggers[name]

	if logger == nil {
		baseLogger := new(BaseLogger)

		baseLogger.name = name
		baseLogger.sinks = config.sinks
		baseLogger.level = config.level

		logger = baseLogger
		loggers[name] = logger
	}

	return logger
}

func (l *BaseLogger) Log(level *LogLevel, message string, data map[string]string) {
	if !l.active(level) {
		return
	}

	for _, sink := range l.sinks {
		record := NewRecord(level, message, data)

		sink.AddRecord(record)
		sink.Flush()
	}
}

func (l *BaseLogger) Fatal(message string) {
	l.Log(LOG_FATAL, message, nil)
}

func (l *BaseLogger) Error(message string) {
	l.Log(LOG_ERROR, message, nil)
}

func (l *BaseLogger) Warn(message string) {
	l.Log(LOG_WARN, message, nil)
}

func (l *BaseLogger) Info(message string) {
	l.Log(LOG_INFO, message, nil)
}

func (l *BaseLogger) Debug(message string) {
	l.Log(LOG_DEBUG, message, nil)
}

func (l *BaseLogger) Debug1(message string) {
	l.Log(LOG_DEBUG1, message, nil)
}

func (l *BaseLogger) Debug2(message string) {
	l.Log(LOG_DEBUG2, message, nil)
}

func (l *BaseLogger) active(level *LogLevel) bool {
	if l.level.priority >= level.priority {
		return true
	}

	return false
}

// For testing
func NumLogger() int {
	return len(loggers)
}
