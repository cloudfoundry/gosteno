package steno

import (
	"encoding/json"
	"fmt"
)

type TaggedLogger struct {
	proxyLogger Logger
	d           map[string]interface{}
}

func NewTaggedLogger(l Logger, d map[string]interface{}) Logger {
	taggedLogger := new(TaggedLogger)

	taggedLogger.proxyLogger = l
	taggedLogger.d = d

	return taggedLogger
}

func (x *TaggedLogger) Log(l LogLevel, m string, d map[string]interface{}) {
	if d != nil {
		e := make(map[string]interface{})

		// Copy the logger's data
		for k, v := range x.d {
			e[k] = v
		}

		// Overwrite specified data
		for k, v := range d {
			e[k] = v
		}

		x.proxyLogger.Log(l, m, e)
	} else {
		x.proxyLogger.Log(l, m, x.d)
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
	data, _ := json.Marshal(x.d)
	proxy, _ := json.Marshal(x.proxyLogger)

	msg := fmt.Sprintf("{\"data\": %s, \"proxy\": %s}", data, proxy)
	return []byte(msg), nil
}
