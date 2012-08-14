package steno

import (
	. "launchpad.net/gocheck"
)

type LoggerSuite struct {
	nSink *nullSink
}

var _ = Suite(&LoggerSuite{})

func (s *LoggerSuite) SetUpTest(c *C) {
	cfg := Config{}
	s.nSink = newNullSink()
	cfg.Sinks = []Sink{s.nSink}
	loggers = make(map[string]*BaseLogger)
	Init(&cfg)
}

func (s *LoggerSuite) TearDownTest(c *C) {
	config = Config{}
	loggers = nil
	loggerRegexp = nil
	loggerRegexpLevel = nil
}

func (s *LoggerSuite) TestLoggersNum(c *C) {
	c.Assert(len(loggers), Equals, 0)
}

func (s *LoggerSuite) TestLoggerLevelActive(c *C) {
	// active is a private method of BaseLogger
	logger := NewLogger("bar").(*BaseLogger)
	logger.level = LOG_INFO
	higherLevels := []*LogLevel{LOG_WARN, LOG_ERROR, LOG_FATAL}
	for _, level := range higherLevels {
		c.Assert(logger.active(level), Equals, true)
	}
}

func (s *LoggerSuite) TestLog(c *C) {
	logger := NewLogger("foobar")
	logger.Info("Hello, world")
	c.Assert(s.nSink.records, HasLen, 1)
	bytes, _ := config.Codec.EncodeRecord(s.nSink.records[0])
	c.Assert(string(bytes), Matches, "{.*Hello, world.*}")
}

func (s *LoggerSuite) TestLevelVisibility(c *C) {
	logger := NewLogger("whatever")
	logger.Info("hello")
	logger.Debug("world")
	// The default level is 'info', so 'debug' will be hidden
	c.Assert(s.nSink.records, HasLen, 1)
}

func (s *LoggerSuite) TestCreatingDupLogger(c *C) {
	logger1 := NewLogger("foobar")
	logger2 := NewLogger("foobar")
	c.Assert(logger1, Equals, logger2)
}
