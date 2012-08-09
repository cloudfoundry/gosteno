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
		logger.level = config.level

		loggers[name] = logger
	}

	return logger
}

func (l *Logger) Log(level *LogLevel, message string) {
	if !l.active(level) {
		return
	}

	for _, sink := range l.sinks {
		record := NewRecord(level, message)

		sink.AddRecord(record)
		sink.Flush()
	}
}

func (l *Logger) Fatal(message string) {
	l.Log(LOG_FATAL, message)
}

func (l *Logger) Error(message string) {
	l.Log(LOG_ERROR, message)
}

func (l *Logger) Warn(message string) {
	l.Log(LOG_WARN, message)
}

func (l *Logger) Info(message string) {
	l.Log(LOG_INFO, message)
}

func (l *Logger) Debug(message string) {
	l.Log(LOG_DEBUG, message)
}

func (l *Logger) Debug1(message string) {
	l.Log(LOG_DEBUG1, message)
}

func (l *Logger) Debug2(message string) {
	l.Log(LOG_DEBUG2, message)
}

func (l *Logger) active(level *LogLevel) bool {
	if l.level.priority >= level.priority {
		return true
	}

	return false
}

// For testing
func NumLogger() int {
	return len(loggers)
}
