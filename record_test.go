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

	c.Assert(record.line, Equals, 15)
	c.Assert(record.method, Equals, "gosteno.(*RecordSuite).TestNewRecordWithLOCEnabled")
	c.Assert(record.file, Matches, ".*record_test.go$")
}

func (s *RecordSuite) TestNewRecordWithLOCDisabled(c *C) {
	message := "Hello, world"
	record := NewRecord(LOG_INFO, message, map[string]string{})

	c.Assert(record.line, Equals, 0)
	c.Assert(record.method, Equals, "")
	c.Assert(record.file, Equals, "")
}
