package steno

import (
	. "launchpad.net/gocheck"
)

type TaggedLoggerSuite struct {
	nSink *nullSink
}

var _ = Suite(&TaggedLoggerSuite{})

func (s *TaggedLoggerSuite) SetUpTest(c *C) {
	s.nSink = newNullSink()

	cfg := Config{}
	cfg.Sinks = []Sink{s.nSink}
	Init(&cfg)
	loggers = make(map[string]*BaseLogger)
}

func (s *TaggedLoggerSuite) TearDownTest(c *C) {
	s.nSink = nil

	config = Config{}
	loggers = nil
	loggerRegexp = nil
	loggerRegexpLevel = nil
}

func (s *TaggedLoggerSuite) TestNewTaggedLogger(c *C) {
	logger := NewLogger("foobar")
	taggedLogger := NewTaggedLogger(logger, map[string]string{"foo": "bar"})
	taggedLogger.Info("Hello")
	taggedLogger.Debug("World")
	// the level of tagged logger should be the same as the derived logger
	c.Assert(s.nSink.records, HasLen, 1)
}

func (s *TaggedLoggerSuite) TestTaggedLogger(c *C) {
	logger := NewLogger("foobar")
	taggedLogger := NewTaggedLogger(logger, map[string]string{"foo": "bar"})
	taggedLogger.Info("Hello")
	bytes, _ := config.Codec.EncodeRecord(s.nSink.records[0])
	c.Assert(string(bytes), Matches, `{.*"foo":"bar".*}`)
}

func (s *TaggedLoggerSuite) TestTaggedLogger2(c *C) {
	logger := NewLogger("whatever")
	taggedLogger := NewTaggedLogger(logger, map[string]string{"foo": "bar"})
	taggedLogger2 := NewTaggedLogger(taggedLogger, map[string]string{"oof": "rab"})
	taggedLogger2.Info("Hello")
	c.Assert(s.nSink.records, HasLen, 1)

	bytes, _ := config.Codec.EncodeRecord(s.nSink.records[0])
	c.Assert(string(bytes), Matches, `{.*"foo":"bar".*}`)
	c.Assert(string(bytes), Matches, `{.*"oof":"rab".*}`)
}
