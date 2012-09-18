package steno

import (
	"log/syslog"
	"sync"
)

type Syslog struct {
	writer *syslog.Writer
	codec  Codec

	sync.Mutex
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
	bytes, _ := s.codec.EncodeRecord(record)
	msg := string(bytes)

	s.Lock()
	defer s.Unlock()

	switch record.Level {
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
		panic("Unknown log level: " + record.Level.name)
	}
}

func (s *Syslog) Flush() {
	// No impl.
}

func (s *Syslog) SetCodec(codec Codec) {
	s.Lock()
	defer s.Unlock()

	s.codec = codec
}

func (s *Syslog) GetCodec() Codec {
	s.Lock()
	defer s.Unlock()

	return s.codec
}

func (s *Syslog) MarshalJSON() ([]byte, error) {
	msg := "{\"type\":\"syslog\"}"
	return []byte(msg), nil
}
