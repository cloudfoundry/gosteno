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

func (x *TaggedLogger) Log(l LogLevel, m string, data map[string]string) {
	if data != nil {
		d := make(map[string]string)

		// data will cover userData if key is the same
		for k, v := range x.data {
			d[k] = v
		}
		for k, v := range data {
			d[k] = v
		}

		x.proxyLogger.Log(l, m, d)
	} else {
		x.proxyLogger.Log(l, m, x.data)
	}
}

func (x *TaggedLogger) Fatal(m string) {
	x.Log(LOG_FATAL, m, nil)
	panic(m)
}

func (x *TaggedLogger) Error(m string) {
	x.Log(LOG_ERROR, m, nil)
}

func (x *TaggedLogger) Warn(m string) {
	x.Log(LOG_WARN, m, nil)
}

func (x *TaggedLogger) Info(m string) {
	x.Log(LOG_INFO, m, nil)
}

func (x *TaggedLogger) Debug(m string) {
	x.Log(LOG_DEBUG, m, nil)
}

func (x *TaggedLogger) Debug1(m string) {
	x.Log(LOG_DEBUG1, m, nil)
}

func (x *TaggedLogger) Debug2(m string) {
	x.Log(LOG_DEBUG2, m, nil)
}

func (x *TaggedLogger) Fatalf(f string, a ...interface{}) {
	x.Fatal(fmt.Sprintf(f, a...))
}

func (x *TaggedLogger) Errorf(f string, a ...interface{}) {
	x.Error(fmt.Sprintf(f, a...))
}

func (x *TaggedLogger) Warnf(f string, a ...interface{}) {
	x.Warn(fmt.Sprintf(f, a...))
}

func (x *TaggedLogger) Infof(f string, a ...interface{}) {
	x.Info(fmt.Sprintf(f, a...))
}

func (x *TaggedLogger) Debugf(f string, a ...interface{}) {
	x.Debug(fmt.Sprintf(f, a...))
}

func (x *TaggedLogger) Debug1f(f string, a ...interface{}) {
	x.Debug1(fmt.Sprintf(f, a...))
}

func (x *TaggedLogger) Debug2f(f string, a ...interface{}) {
	x.Debug2(fmt.Sprintf(f, a...))
}

func (x *TaggedLogger) MarshalJSON() ([]byte, error) {
	data, _ := json.Marshal(x.data)
	proxy, _ := json.Marshal(x.proxyLogger)

	msg := fmt.Sprintf("{\"data\": %s, \"proxy\": %s}", data, proxy)
	return []byte(msg), nil
}
