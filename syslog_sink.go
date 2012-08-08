package steno

import (
	"log/syslog"
	"fmt"
)

// FIXME: Fill the full map
var levelMap = map[string]syslog.Priority{
	"info":  syslog.LOG_INFO,
	"debug": syslog.LOG_DEBUG,
}

type Syslog struct {
	writer *syslog.Writer
}

func NewSyslogSink() *Syslog {
	writer, err := syslog.New(syslog.LOG_DEBUG, "")
	if err != nil {
		panic(err)
	}

	syslog := new(Syslog)
	syslog.writer = writer
	return syslog
}

func (syslog *Syslog) AddRecord(record *Record) {
	// TODO: unified format
	msg := fmt.Sprintf("[%s] %s %s\n", record.timestamp, record.level.name, record.message)
	// FIXME: use info defaultly
	syslog.writer.Info(msg)
}

func (syslog *Syslog) Flush() {
	// No impl.
}
