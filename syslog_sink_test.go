package steno

import (
	. "launchpad.net/gocheck"
)

type SyslogSinkSuite struct {
	syslogSink *Syslog
}

var _ = Suite(&SyslogSinkSuite{})

func (s *SyslogSinkSuite) SetUpTest(c *C) {
	s.syslogSink = NewSyslogSink()
}

func (s *SyslogSinkSuite) TearDownTest(c *C) {
	s.syslogSink = nil
}

func (s *SyslogSinkSuite) TestSyslogMarshalJson(c *C) {
	bytes, _ := s.syslogSink.MarshalJSON()
	c.Assert(string(bytes), Equals, `{"type":"syslog"}`)
}
