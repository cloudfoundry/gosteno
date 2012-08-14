package steno

import (
	. "launchpad.net/gocheck"
)

type TaggedLoggerSuit struct {
}

var _ = Suite(&TaggedLoggerSuit{})

func (s *TaggedLoggerSuit) SetUpTest(c *C) {

}

func (s *TaggedLoggerSuit) TearDownTest(c *C) {
}

func (s *TaggedLoggerSuit) TestTaggedLogger(c *C) {
	nSink := new(NullSink)
	nSink.records = make([]*Record, 0, 10)
	sinks := []Sink{nSink}
	baseLogger := new(BaseLogger)
	baseLogger.name = "foobar"
	baseLogger.sinks = sinks
	baseLogger.level = LOG_INFO

	data := make(map[string]string)
	data["foo"] = "bar"
	taggedLogger := NewTaggedLogger(baseLogger, data)
	taggedLogger.Info("Hello")
	c.Assert(nSink.records, HasLen, 1)
}
