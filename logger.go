package steno

import (
	"encoding/json"
	"fmt"
	"log"
)

type L interface {
	Level() LogLevel
	Log(x LogLevel, m string, d map[string]interface{})
}

type Logger struct {
	L
}

type BaseLogger struct {
	name  string
	sinks []Sink
	level LogLevel
}

func (l *BaseLogger) Level() LogLevel {
	return l.level
}

func (l *BaseLogger) Log(x LogLevel, m string, d map[string]interface{}) {
	if l.Level().Priority < x.Priority {
		return
	}

	r := NewRecord(l.name, x, m, d)
	for _, s := range l.sinks {
		s.AddRecord(r)
		s.Flush()
	}

	if x == LOG_FATAL {
		panic(m)
	}
}

func (l *BaseLogger) MarshalJSON() ([]byte, error) {
	sinks := "["
	for i, sink := range l.sinks {
		m, err := json.Marshal(sink)
		if err != nil {
			log.Println(err)
		}
		sinks += string(m)
		if i != len(l.sinks)-1 {
			sinks += ","
		}
	}
	sinks += "]"
	msg := fmt.Sprintf("{\"level\": \"%s\", \"sinks\": %s}", l.level.Name, sinks)
	return []byte(msg), nil
}

func (l Logger) Fatal(m string) {
	l.Log(LOG_FATAL, m, nil)
}

func (l Logger) Error(m string) {
	l.Log(LOG_ERROR, m, nil)
}

func (l Logger) Warn(m string) {
	l.Log(LOG_WARN, m, nil)
}

func (l Logger) Info(m string) {
	l.Log(LOG_INFO, m, nil)
}

func (l Logger) Debug(m string) {
	l.Log(LOG_DEBUG, m, nil)
}

func (l Logger) Debug1(m string) {
	l.Log(LOG_DEBUG1, m, nil)
}

func (l Logger) Debug2(m string) {
	l.Log(LOG_DEBUG2, m, nil)
}

func (l Logger) Fatald(d map[string]interface{}, m string) {
	l.Log(LOG_FATAL, m, d)
}

func (l Logger) Errord(d map[string]interface{}, m string) {
	l.Log(LOG_ERROR, m, d)
}

func (l Logger) Warnd(d map[string]interface{}, m string) {
	l.Log(LOG_WARN, m, d)
}

func (l Logger) Infod(d map[string]interface{}, m string) {
	l.Log(LOG_INFO, m, d)
}

func (l Logger) Debugd(d map[string]interface{}, m string) {
	l.Log(LOG_DEBUG, m, d)
}

func (l Logger) Debug1d(d map[string]interface{}, m string) {
	l.Log(LOG_DEBUG1, m, d)
}

func (l Logger) Debug2d(d map[string]interface{}, m string) {
	l.Log(LOG_DEBUG2, m, d)
}

func (l Logger) Fatalf(f string, a ...interface{}) {
	l.Log(LOG_FATAL, fmt.Sprintf(f, a...), nil)
}

func (l Logger) Errorf(f string, a ...interface{}) {
	l.Log(LOG_ERROR, fmt.Sprintf(f, a...), nil)
}

func (l Logger) Warnf(f string, a ...interface{}) {
	l.Log(LOG_WARN, fmt.Sprintf(f, a...), nil)
}

func (l Logger) Infof(f string, a ...interface{}) {
	l.Log(LOG_INFO, fmt.Sprintf(f, a...), nil)
}

func (l Logger) Debugf(f string, a ...interface{}) {
	l.Log(LOG_DEBUG, fmt.Sprintf(f, a...), nil)
}

func (l Logger) Debug1f(f string, a ...interface{}) {
	l.Log(LOG_DEBUG1, fmt.Sprintf(f, a...), nil)
}

func (l Logger) Debug2f(f string, a ...interface{}) {
	l.Log(LOG_DEBUG2, fmt.Sprintf(f, a...), nil)
}

func (l Logger) Fataldf(d map[string]interface{}, f string, a ...interface{}) {
	l.Log(LOG_FATAL, fmt.Sprintf(f, a...), d)
}

func (l Logger) Errordf(d map[string]interface{}, f string, a ...interface{}) {
	l.Log(LOG_ERROR, fmt.Sprintf(f, a...), d)
}

func (l Logger) Warndf(d map[string]interface{}, f string, a ...interface{}) {
	l.Log(LOG_WARN, fmt.Sprintf(f, a...), d)
}

func (l Logger) Infodf(d map[string]interface{}, f string, a ...interface{}) {
	l.Log(LOG_INFO, fmt.Sprintf(f, a...), d)
}

func (l Logger) Debugdf(d map[string]interface{}, f string, a ...interface{}) {
	l.Log(LOG_DEBUG, fmt.Sprintf(f, a...), d)
}

func (l Logger) Debug1df(d map[string]interface{}, f string, a ...interface{}) {
	l.Log(LOG_DEBUG1, fmt.Sprintf(f, a...), d)
}

func (l Logger) Debug2df(d map[string]interface{}, f string, a ...interface{}) {
	l.Log(LOG_DEBUG2, fmt.Sprintf(f, a...), d)
}
