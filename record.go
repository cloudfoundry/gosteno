package steno

import (
	"time"
)

// FIXME: Missing fields
type Record struct {
	timestamp time.Time
	message string
	level *LogLevel
}

func NewRecord(level *LogLevel, message string) *Record {
	record := new(Record)

	record.timestamp = time.Now()
	record.message = message
	record.level = level

	return record
}
