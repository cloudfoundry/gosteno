package steno

import (
	"encoding/json"
	"fmt"
)

type TaggedLogger struct {
	proxyLogger Logger
	data        map[string]string
}

// tagged logger doesn't have name, so far
func NewTaggedLogger(logger Logger, data map[string]string) Logger {
	taggedLogger := new(TaggedLogger)

	taggedLogger.proxyLogger = logger
	taggedLogger.data = data

	return taggedLogger
}

func (l *TaggedLogger) Log(level *LogLevel, message string, data map[string]string) {
	if data != nil {
		d := make(map[string]string)

		// data will cover userData if key is the same
		for k, v := range l.data {
			d[k] = v
		}
		for k, v := range data {
			d[k] = v
		}

		l.proxyLogger.Log(level, message, d)
	} else {
		l.proxyLogger.Log(level, message, l.data)
	}
}

func (l *TaggedLogger) Fatal(message string) {
	l.Log(LOG_FATAL, message, nil)
}

func (l *TaggedLogger) Error(message string) {
	l.Log(LOG_ERROR, message, nil)
}

func (l *TaggedLogger) Warn(message string) {
	l.Log(LOG_WARN, message, nil)
}

func (l *TaggedLogger) Info(message string) {
	l.Log(LOG_INFO, message, nil)
}

func (l *TaggedLogger) Debug(message string) {
	l.Log(LOG_DEBUG, message, nil)
}

func (l *TaggedLogger) Debug1(message string) {
	l.Log(LOG_DEBUG1, message, nil)
}

func (l *TaggedLogger) Debug2(message string) {
	l.Log(LOG_DEBUG2, message, nil)
}

func (l *TaggedLogger) Fatalf(format string, a ...interface{}) {
	l.Fatal(fmt.Sprintf(format, a...))
}

func (l *TaggedLogger) Errorf(format string, a ...interface{}) {
	l.Error(fmt.Sprintf(format, a...))
}

func (l *TaggedLogger) Warnf(format string, a ...interface{}) {
	l.Warn(fmt.Sprintf(format, a...))
}

func (l *TaggedLogger) Infof(format string, a ...interface{}) {
	l.Info(fmt.Sprintf(format, a...))
}

func (l *TaggedLogger) Debugf(format string, a ...interface{}) {
	l.Debug(fmt.Sprintf(format, a...))
}

func (l *TaggedLogger) Debug1f(format string, a ...interface{}) {
	l.Debug1(fmt.Sprintf(format, a...))
}

func (l *TaggedLogger) Debug2f(format string, a ...interface{}) {
	l.Debug2(fmt.Sprintf(format, a...))
}

func (l *TaggedLogger) MarshalJSON() ([]byte, error) {
	data, _ := json.Marshal(l.data)
	proxy, _ := json.Marshal(l.proxyLogger)

	msg := fmt.Sprintf("{\"data\": %s, \"proxy\": %s}", data, proxy)
	return []byte(msg), nil
}
