package steno

import (
	"encoding/json"
	. "launchpad.net/gocheck"
)

type LogLevelSuite struct {
}

var _ = Suite(&LogLevelSuite{})

func (s *LogLevelSuite) TestNewLogLevel(c *C) {
	level := newLogLevel("foobar", 100)
	c.Assert(level, NotNil)
	c.Assert(level.name, Equals, "foobar")
	c.Assert(level.priority, Equals, 100)
}

func (s *LogLevelSuite) TestGetLevel(c *C) {
	infoLevel, err := GetLogLevel("info")
	c.Assert(infoLevel, Equals, LOG_INFO)
	c.Assert(err, IsNil)
}

func (s *LogLevelSuite) TestGetNotExistLevel(c *C) {
	notExistLevel, err := GetLogLevel("foobar")
	c.Assert(notExistLevel, IsNil)
	c.Assert(err, NotNil)
}

func (s *LogLevelSuite) TestImplementedInterfaces(c *C) {
	var v interface{} = newLogLevel("foobar", 1)

	// *LogLevel should implement json.Marshaler interface
	_, ok := v.(json.Marshaler)
	c.Assert(ok, Equals, true)

	// *LogLevel should implement json.Unmarshaler interface
	_, ok = v.(json.Unmarshaler)
	c.Assert(ok, Equals, true)
}
