package steno

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
