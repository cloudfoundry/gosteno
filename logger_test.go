package steno

import (
	. "launchpad.net/gocheck"
	"os"
)

type LoggerSuite struct {
}

var _ = Suite(&LoggerSuite{})

func (s *LoggerSuite) SetUpTest(c *C) {
	sinks := []Sink{NewIOSink(os.Stdout)}
	cfg := Config{}
	cfg.Sinks = sinks
	cfg.Codec = NewJsonCodec()
	cfg.Level = LOG_INFO
	loggers = make(map[string]*BaseLogger)
	Init(&cfg)
}

func (s *LoggerSuite) TearDownTest(c *C) {
	config = Config{}
	loggers = nil
}

func (s *LoggerSuite) TestLoggersNum(c *C) {
	c.Assert(len(loggers), Equals, 0)
}

func (s *LoggerSuite) TestDefaultLevel(c *C) {

}

func (s *LoggerSuite) TestLoggerName(c *C) {
	logger := NewLogger("test").(*BaseLogger)
	c.Assert(logger.name, Equals, "test")
}

func (s *LoggerSuite) TestLoggerLevelActive(c *C) {
	logger := NewLogger("foo").(*BaseLogger)
	defaultLevel := LOG_INFO
	c.Assert(logger.active(defaultLevel), Equals, true)
}

func (s *LoggerSuite) TestLoggerLevelActive2(c *C) {
	logger := NewLogger("bar").(*BaseLogger)
	logger.level = LOG_INFO
	higherLevels := []*LogLevel{LOG_WARN, LOG_ERROR, LOG_ERROR}
	for _, level := range higherLevels {
		c.Assert(logger.active(level), Equals, true)
	}
}

func (s *LoggerSuite) TestCreatingDupLogger(c *C) {
	logger1 := NewLogger("foobar")
	logger2 := NewLogger("foobar")
	c.Assert(logger1, Equals, logger2)
}
