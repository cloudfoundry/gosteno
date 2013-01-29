package steno

import (
	"os"
	"runtime"
	"strings"
	"time"
)

// FIXME: Missing fields
type Record struct {
	Timestamp float64
	Pid       int
	Source    string
	Level     LogLevel
	Message   string
	Data      map[string]string
	File      string
	Line      int
	Method    string
}

var pid int

func init() {
	pid = os.Getpid()
}

func NewRecord(s string, l LogLevel, m string, d map[string]string) *Record {
	r := &Record{
		Timestamp: float64(time.Now().UnixNano()) / 1000000000,
		Pid:       pid,
		Source:    s,
		Level:     l,
		Message:   m,
		Data:      d,
	}

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
		r.File = file
		r.Line = line
		r.Method = function.Name()
	}

	return r
}
