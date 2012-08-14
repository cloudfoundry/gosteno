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
