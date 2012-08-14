package steno

import (
	. "launchpad.net/gocheck"
)

type TaggedLoggerSuite struct {
}

var _ = Suite(&TaggedLoggerSuite{})

func (s *TaggedLoggerSuite) TestTaggedLogger(c *C) {
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
