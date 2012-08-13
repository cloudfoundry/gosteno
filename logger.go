package steno

import "fmt"

type Logger interface {
	Log(level *LogLevel, message string, data map[string]string)
	Fatal(message string)
	Error(message string)
	Warn(message string)
	Info(message string)
	Debug(message string)
	Debug1(message string)
	Debug2(message string)

	Fatalf(format string, a ...interface{})
	Errorf(format string, a ...interface{})
	Warnf(format string, a ...interface{})
	Infof(format string, a ...interface{})
	Debugf(format string, a ...interface{})
	Debug1f(format string, a ...interface{})
	Debug2f(format string, a ...interface{})
}

type BaseLogger struct {
	name string
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

func (l *BaseLogger) Fatalf(format string, a ...interface{}) {
	l.Fatal(fmt.Sprintf(format, a...))
}

func (l *BaseLogger) Errorf(format string, a ...interface{}) {
	l.Error(fmt.Sprintf(format, a...))
}

func (l *BaseLogger) Warnf(format string, a ...interface{}) {
	l.Warn(fmt.Sprintf(format, a...))
}

func (l *BaseLogger) Infof(format string, a ...interface{}) {
	l.Info(fmt.Sprintf(format, a...))
}

func (l *BaseLogger) Debugf(format string, a ...interface{}) {
	l.Debug(fmt.Sprintf(format, a...))
}

func (l *BaseLogger) Debug1f(format string, a ...interface{}) {
	l.Debug1(fmt.Sprintf(format, a...))
}

func (l *BaseLogger) Debug2f(format string, a ...interface{}) {
	l.Debug2(fmt.Sprintf(format, a...))
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
