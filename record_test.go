package steno

import (
	. "launchpad.net/gocheck"
)

type RecordSuite struct {
}

var _ = Suite(&RecordSuite{})

func (s *RecordSuite) TestNewRecordWithLOCEnabled(c *C) {
	message := "Hello, GOSTENO"
	config.EnableLOC = true
	record := NewRecord(LOG_INFO, message, map[string]string{}) // Line 15
	config.EnableLOC = false

	c.Assert(record.Line, Equals, 15)
	c.Assert(record.Method, Matches, `.*\.\(\*RecordSuite\)\.TestNewRecordWithLOCEnabled`)
	c.Assert(record.File, Matches, ".*record_test.go$")
}

func (s *RecordSuite) TestNewRecordWithLOCDisabled(c *C) {
	message := "Hello, world"
	record := NewRecord(LOG_INFO, message, map[string]string{})

	c.Assert(record.Line, Equals, 0)
	c.Assert(record.Method, Equals, "")
	c.Assert(record.File, Equals, "")
}
