package steno

import (
	"os"
	"runtime"
	"strings"
	"time"
)

type Record struct {
	Timestamp float64           `json:"timestamp"`
	Pid       int               `json:"process_id"`
	Source    string            `json:"source"`
	Level     LogLevel          `json:"log_level"`
	Message   string            `json:"message"`
	Data      map[string]string `json:"data"`
	File      string            `json:"file"`
	Line      int               `json:"line"`
	Method    string            `json:"method"`
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
