package steno

import (
	"encoding/json"
	"fmt"
)

type LogLevel struct {
	name     string
	priority int
}

var (
	LOG_OFF    = newLogLevel("off", 0)
	LOG_FATAL  = newLogLevel("fatal", 1)
	LOG_ERROR  = newLogLevel("error", 5)
	LOG_WARN   = newLogLevel("warn", 10)
	LOG_INFO   = newLogLevel("info", 15)
	LOG_DEBUG  = newLogLevel("debug", 16)
	LOG_DEBUG1 = newLogLevel("debug1", 17)
	LOG_DEBUG2 = newLogLevel("debug2", 18)
	LOG_ALL    = newLogLevel("all", 30)
)

var levels = map[string]*LogLevel{}

func newLogLevel(name string, priority int) *LogLevel {
	level := new(LogLevel)

	level.name = name
	level.priority = priority

	levels[name] = level

	return level
}

func GetLogLevel(name string) (*LogLevel, error) {
	if level, ok := levels[name]; ok {
		return level, nil
	}
	err := fmt.Errorf("No level with that name exists : %s", name)
	return nil, err
}

func (level *LogLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(level.name)
}

func (level *LogLevel) UnmarshalJSON(data []byte) error {
	var n string
	err := json.Unmarshal(data, &n)
	if err != nil {
		return err
	}
	if l, err := GetLogLevel(n); err != nil {
		return err
	} else {
		*level = *l
	}

	return nil
}

func (level *LogLevel) String() string {
	return level.name
}
