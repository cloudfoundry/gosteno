package steno

import (
	. "launchpad.net/gocheck"
	"os"
	"testing"
)

func Test(t *testing.T) {
	TestingT(t)
}

type StenoSuite struct {
}

var _ = Suite(&StenoSuite{})

func (s *StenoSuite) SetUpTest(c *C) {
	cfg := Config{}
	cfg.Sinks = []Sink{NewIOSink(os.Stdout)}
	Init(&cfg)
	loggers = make(map[string]*BaseLogger)
}

func (s *StenoSuite) TearDownTest(c *C) {
	config = Config{}
	loggers = nil
	loggerRegexp = nil
	loggerRegexpLevel = nil
}

func (s *StenoSuite) TestInitLoggers(c *C) {
	c.Assert(loggers, HasLen, 0)
}

func (s *StenoSuite) TestDefaultConfig(c *C) {
	c.Assert(config.Codec, Equals, JSON_CODEC)
	c.Assert(config.Level, Equals, LOG_INFO)
	c.Assert(config.Port, Equals, 0)
}

func (s *StenoSuite) TestLoggersInJson(c *C) {
	c.Assert(loggersInJson(), Matches, "{.*}")
}

func (s *StenoSuite) TestSetLoggerRegexp(c *C) {
	// level is a field of BaseLogger, hence type cast is needed
	logger1 := NewLogger("test").(*BaseLogger)
	logger2 := NewLogger("test2").(*BaseLogger)
	logger3 := NewLogger("test3").(*BaseLogger)

	c.Assert(logger1.level, Equals, LOG_INFO)
	c.Assert(logger2.level, Equals, LOG_INFO)
	c.Assert(logger3.level, Equals, LOG_INFO)

	SetLoggerRegexp("te", LOG_FATAL)
	c.Assert(logger1.level, Equals, LOG_FATAL)
	c.Assert(logger2.level, Equals, LOG_FATAL)
	c.Assert(logger3.level, Equals, LOG_FATAL)

	SetLoggerRegexp("test", LOG_ERROR)
	c.Assert(logger1.level, Equals, LOG_ERROR)
	c.Assert(logger2.level, Equals, LOG_ERROR)
	c.Assert(logger3.level, Equals, LOG_ERROR)

	SetLoggerRegexp("test$", LOG_WARN)
	c.Assert(logger1.level, Equals, LOG_WARN)
	c.Assert(logger2.level, Equals, LOG_INFO)
	c.Assert(logger3.level, Equals, LOG_INFO)
}

func (s *StenoSuite) TestClearLoggerRegexp(c *C) {
	// level is a field of BaseLogger, hence type cast is needed
	logger1 := NewLogger("test").(*BaseLogger)
	logger2 := NewLogger("test2").(*BaseLogger)
	logger3 := NewLogger("test3").(*BaseLogger)

	c.Assert(logger1.level, Equals, LOG_INFO)
	c.Assert(logger2.level, Equals, LOG_INFO)
	c.Assert(logger3.level, Equals, LOG_INFO)

	SetLoggerRegexp("test", LOG_FATAL)
	c.Assert(logger1.level, Equals, LOG_FATAL)
	c.Assert(logger2.level, Equals, LOG_FATAL)
	c.Assert(logger3.level, Equals, LOG_FATAL)

	ClearLoggerRegexp()
	c.Assert(logger1.level, Equals, LOG_INFO)
	c.Assert(logger2.level, Equals, LOG_INFO)
	c.Assert(logger3.level, Equals, LOG_INFO)
}
