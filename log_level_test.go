package steno

import (
	"encoding/json"
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

func (s *LogLevelSuite) TestImplementedInterfaces(c *C) {
	var v interface{} = NewLogLevel("foobar", 1)

	// *LogLevel should implement json.Marshaler interface
	_, ok := v.(json.Marshaler)
	c.Assert(ok, Equals, true)

	// *LogLevel should implement json.Unmarshaler interface
	_, ok = v.(json.Unmarshaler)
	c.Assert(ok, Equals, true)
}
