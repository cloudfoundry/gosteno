package steno

import (
	"log/syslog"
)

// FIXME: Fill the full map
var levelMap = map[string]syslog.Priority{
	"info":  syslog.LOG_INFO,
	"debug": syslog.LOG_DEBUG,
}

type Syslog struct {
	writer *syslog.Writer
	codec Codec
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

func (s *Syslog) AddRecord(record *Record) {
	msg := s.codec.EncodeRecord(record)

	// FIXME: use info defaultly
	s.writer.Info(msg)
}

func (s *Syslog) Flush() {
	// No impl.
}

func (s *Syslog) SetCodec(codec Codec) {
	s.codec = codec
}

func (s *Syslog) GetCodec() Codec {
	return s.codec
}
