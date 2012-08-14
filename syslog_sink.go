package steno

import (
	"log/syslog"
)

type Syslog struct {
	writer *syslog.Writer
	codec  Codec
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

	switch record.level {
	case LOG_FATAL:
		s.writer.Crit(msg)
	case LOG_ERROR:
		s.writer.Err(msg)
	case LOG_WARN:
		s.writer.Warning(msg)
	case LOG_INFO:
		s.writer.Info(msg)
	case LOG_DEBUG, LOG_DEBUG1, LOG_DEBUG2:
		s.writer.Debug(msg)
	default:
		panic("Unknown log level: " + record.level.name)
	}
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
