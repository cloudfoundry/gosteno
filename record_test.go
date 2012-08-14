package steno

import (
	. "launchpad.net/gocheck"
)

type RecordSuite struct {
}

var _ = Suite(&RecordSuite{})

func (s *RecordSuite) TestNewRecord(c *C) {
	message := "Hello, GOSTENO"
	data := make(map[string]string)
	record := NewRecord(LOG_INFO, message, data)
	c.Assert(record, NotNil)
}
