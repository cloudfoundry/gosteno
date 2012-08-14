package steno

import (
	"time"
	"runtime"
	"strings"
)

// FIXME: Missing fields
type Record struct {
	timestamp time.Time
	message   string
	level     *LogLevel
	data      map[string]string
	file      string
	method    string
	line      int
}

func NewRecord(level *LogLevel, message string, data map[string]string) *Record {
	record := new(Record)

	record.timestamp = time.Now()
	record.message = message
	record.level = level
	record.data = data

	var function *runtime.Func
	var file string
	var line int

	pc := make([]uintptr, 50)
	nptrs := runtime.Callers(2, pc)
	for i := 0; i < nptrs; i++ {
		function = runtime.FuncForPC(pc[i])
		file, line = function.FileLine(pc[i])
		if !strings.HasSuffix(file, "logger.go") {
			break
		}
	}
	record.file = file
	record.method = function.Name()
	record.line = line

	return record
}
