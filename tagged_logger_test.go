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
	nSink.codec = NewJsonCodec()
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

	msgBytes, _ := nSink.codec.EncodeRecord(nSink.records[0])
	c.Assert((string)(msgBytes), Matches, `{.*"foo":"bar".*}`)
}

func (s TaggedLoggerSuit) TestTag(c *C) {
	baseLogger := new(BaseLogger)
	baseLogger.name = "foobar"
	baseLogger.level = LOG_INFO

	data := make(map[string]string)
	data["foo"] = "bar"
	logger := NewTaggedLogger(baseLogger, data)
	taggedLogger := logger.(*TaggedLogger)

	newData := make(map[string]string)
	newData["bar"] = "foo"
	newLogger := taggedLogger.Tag(newData)
	newTagLogger := newLogger.(*TaggedLogger)
	c.Assert(taggedLogger.data["bar"], Equals, "")
	c.Assert(newTagLogger.data["foo"], Equals, "bar")
	c.Assert(newTagLogger.data["bar"], Equals, "foo")
}
