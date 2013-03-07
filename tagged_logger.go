package steno

import (
	"encoding/json"
	"fmt"
)

type TaggedLogger struct {
	Logger

	d map[string]interface{}
}

func NewTaggedLogger(l Logger, d map[string]interface{}) Logger {
	tl := &TaggedLogger{
		Logger: l,
		d:      d,
	}

	return Logger{tl}
}

func (l *TaggedLogger) Log(x LogLevel, m string, d map[string]interface{}) {
	if d != nil {
		e := make(map[string]interface{})

		// Copy the logger's data
		for k, v := range l.d {
			e[k] = v
		}

		// Overwrite specified data
		for k, v := range d {
			e[k] = v
		}

		l.Logger.Log(x, m, e)
	} else {
		l.Logger.Log(x, m, l.d)
	}
}

func (l *TaggedLogger) MarshalJSON() ([]byte, error) {
	data, _ := json.Marshal(l.d)
	proxy, _ := json.Marshal(l.Logger)

	msg := fmt.Sprintf("{\"data\": %s, \"proxy\": %s}", data, proxy)
	return []byte(msg), nil
}
