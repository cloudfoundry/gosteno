package steno

import (
	"encoding/json"
	"fmt"
	"log"
)

type Logger interface {
	json.Marshaler

	Log(level LogLevel, m string, data map[string]string)
	Fatal(m string)
	Error(m string)
	Warn(m string)
	Info(m string)
	Debug(m string)
	Debug1(m string)
	Debug2(m string)

	Fatalf(f string, a ...interface{})
	Errorf(f string, a ...interface{})
	Warnf(f string, a ...interface{})
	Infof(f string, a ...interface{})
	Debugf(f string, a ...interface{})
	Debug1f(f string, a ...interface{})
	Debug2f(f string, a ...interface{})
}

type BaseLogger struct {
	name  string
	sinks []Sink
	level LogLevel
}

func (x *BaseLogger) Log(l LogLevel, m string, data map[string]string) {
	if !x.active(l) {
		return
	}

	record := NewRecord(x.name, l, m, data)

	for _, sink := range x.sinks {
		sink.AddRecord(record)
		sink.Flush()
	}
}

func (x *BaseLogger) Fatal(m string) {
	x.Log(LOG_FATAL, m, nil)
	panic(m)
}

func (x *BaseLogger) Error(m string) {
	x.Log(LOG_ERROR, m, nil)
}

func (x *BaseLogger) Warn(m string) {
	x.Log(LOG_WARN, m, nil)
}

func (x *BaseLogger) Info(m string) {
	x.Log(LOG_INFO, m, nil)
}

func (x *BaseLogger) Debug(m string) {
	x.Log(LOG_DEBUG, m, nil)
}

func (x *BaseLogger) Debug1(m string) {
	x.Log(LOG_DEBUG1, m, nil)
}

func (x *BaseLogger) Debug2(m string) {
	x.Log(LOG_DEBUG2, m, nil)
}

func (x *BaseLogger) Fatalf(f string, a ...interface{}) {
	x.Fatal(fmt.Sprintf(f, a...))
}

func (x *BaseLogger) Errorf(f string, a ...interface{}) {
	x.Error(fmt.Sprintf(f, a...))
}

func (x *BaseLogger) Warnf(f string, a ...interface{}) {
	x.Warn(fmt.Sprintf(f, a...))
}

func (x *BaseLogger) Infof(f string, a ...interface{}) {
	x.Info(fmt.Sprintf(f, a...))
}

func (x *BaseLogger) Debugf(f string, a ...interface{}) {
	x.Debug(fmt.Sprintf(f, a...))
}

func (x *BaseLogger) Debug1f(f string, a ...interface{}) {
	x.Debug1(fmt.Sprintf(f, a...))
}

func (x *BaseLogger) Debug2f(f string, a ...interface{}) {
	x.Debug2(fmt.Sprintf(f, a...))
}

func (x *BaseLogger) MarshalJSON() ([]byte, error) {
	sinks := "["
	for i, sink := range x.sinks {
		m, err := json.Marshal(sink)
		if err != nil {
			log.Println(err)
		}
		sinks += string(m)
		if i != len(x.sinks)-1 {
			sinks += ","
		}
	}
	sinks += "]"
	msg := fmt.Sprintf("{\"level\": \"%s\", \"sinks\": %s}", x.level.Name, sinks)
	return []byte(msg), nil
}

func (x *BaseLogger) active(level LogLevel) bool {
	if x.level.Priority >= level.Priority {
		return true
	}

	return false
}

// For testing
func NumLogger() int {
	return len(loggers)
}
