package steno

import (
	. "launchpad.net/gocheck"
)

type LoggerSuite struct {
	nSink *nullSink
}

var _ = Suite(&LoggerSuite{})

func NumLogger() int {
	return len(loggers)
}

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

func (s *LoggerSuite) TestLogLevel(c *C) {
	bl := &BaseLogger{
		name:  "bar",
		level: LOG_INFO,
		sinks: []Sink{&nullSink{}},
	}

	higher := []LogLevel{LOG_INFO, LOG_WARN, LOG_ERROR}
	for _, l := range higher {
		s := &nullSink{}
		bl.sinks = []Sink{s}
		bl.Log(l, "hello", nil)
		c.Assert(len(s.records), Equals, 1)
	}

	lower := []LogLevel{LOG_DEBUG, LOG_DEBUG1, LOG_DEBUG2, LOG_ALL}
	for _, l := range lower {
		s := &nullSink{}
		bl.sinks = []Sink{s}
		bl.Log(l, "hello", nil)
		c.Assert(len(s.records), Equals, 0)
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

func (s *LoggerSuite) TestPanic(c *C) {
	logger := NewLogger("foobar")
	c.Assert(func() { logger.Fatal("fail!") }, PanicMatches, "fail!")
	c.Assert(func() { logger.Fatalf("fail!%s", "fail!") }, PanicMatches, "fail!fail!")

	t := NewTaggedLogger(logger, map[string]interface{}{"foo": "bar"})
	c.Assert(func() { t.Fatal("panic") }, PanicMatches, "panic")
	c.Assert(func() { t.Fatalf("panic!%s", "panic!") }, PanicMatches, "panic!panic!")
}
