package steno

import (
	"runtime"
	"strings"
	"time"
)

// FIXME: Missing fields
type Record struct {
	Timestamp float64
	Message   string
	Level     *LogLevel
	Data      map[string]string
	File      string
	Method    string
	Line      int
}

func NewRecord(level *LogLevel, message string, data map[string]string) *Record {
	record := new(Record)

	record.Timestamp = float64(time.Now().UnixNano()) / 1000000000
	record.Message = message
	record.Level = level
	record.Data = data

	if config.EnableLOC {
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
		record.File = file
		record.Method = function.Name()
		record.Line = line
	}

	return record
}
