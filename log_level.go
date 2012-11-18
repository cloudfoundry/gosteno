package steno

import (
	"encoding/json"
	"fmt"
)

type LogLevel struct {
	Name     string
	Priority int
}

var (
	LOG_OFF    = defineLogLevel("off", 0)
	LOG_FATAL  = defineLogLevel("fatal", 1)
	LOG_ERROR  = defineLogLevel("error", 5)
	LOG_WARN   = defineLogLevel("warn", 10)
	LOG_INFO   = defineLogLevel("info", 15)
	LOG_DEBUG  = defineLogLevel("debug", 16)
	LOG_DEBUG1 = defineLogLevel("debug1", 17)
	LOG_DEBUG2 = defineLogLevel("debug2", 18)
	LOG_ALL    = defineLogLevel("all", 30)
)

var levels = make(map[string]LogLevel)

func defineLogLevel(n string, p int) LogLevel {
	level := LogLevel{Name: n, Priority: p}

	levels[n] = level

	return level
}

func GetLogLevel(name string) (LogLevel, error) {
	var level LogLevel

	if level, ok := levels[name]; ok {
		return level, nil
	}

	err := fmt.Errorf("Undefined log level: %s", name)
	return level, err
}

func (level LogLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(level.Name)
}

func (level *LogLevel) UnmarshalJSON(data []byte) error {
	var n string

	err := json.Unmarshal(data, &n)
	if err != nil {
		return err
	}

	l, err := GetLogLevel(n)
	if err != nil {
		return err
	}

	*level = l

	return nil
}

func (level LogLevel) String() string {
	return level.Name
}
