package steno

import (
	"time"
)

// FIXME: Missing fields
type Record struct {
	timestamp time.Time
	message   string
	level     *LogLevel
	data      map[string]string
}

func NewRecord(level *LogLevel, message string, data map[string]string) *Record {
	record := new(Record)

	record.timestamp = time.Now()
	record.message = message
	record.level = level
	record.data = data

	return record
}
