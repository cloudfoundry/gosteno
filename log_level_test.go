package steno

import (
	. "launchpad.net/gocheck"
)

type LogLevelSuite struct {
}

var _ = Suite(&LogLevelSuite{})

func (s *LogLevelSuite) TestNewLogLevel(c *C) {
	level := NewLogLevel("foobar", 100)
	c.Assert(level, NotNil)
	c.Assert(level.name, Equals, "foobar")
	c.Assert(level.priority, Equals, 100)
}

func (s *LogLevelSuite) TestLookupLevel(c *C) {
	infoLevel := lookupLevel("info")
	c.Assert(infoLevel, Equals, LOG_INFO)
}

func (s *LogLevelSuite) TestLookupNotExistLevel(c *C) {
	notExistLevel := lookupLevel("foobar")
	c.Assert(notExistLevel, IsNil)
}
