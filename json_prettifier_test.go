package steno

import (
	. "launchpad.net/gocheck"
)

type JsonPrettifierSuite struct {
}

var _ = Suite(&JsonPrettifierSuite{})

func (s *JsonPrettifierSuite) TestConst(c *C) {
	c.Assert(EXCLUDE_NONE, Equals, 0)
	c.Assert(EXCLUDE_TIMESTAMP, Equals, 1)
	c.Assert(EXCLUDE_LEVEL, Equals, 1<<5)
	c.Assert(EXCLUDE_LINE, Equals, 1<<3)
}

func (s *JsonPrettifierSuite) TestConstOrder(c *C) {
	record := NewRecord(LOG_INFO, "Hello, world", map[string]string{"foo": "bar"})

	prettifier1 := NewJsonPrettifier(EXCLUDE_FILE | EXCLUDE_DATA)
	bytes1, _ := prettifier1.PrettifyEntry(record)

	prettifier2 := NewJsonPrettifier(EXCLUDE_DATA | EXCLUDE_FILE)
	bytes2, _ := prettifier2.PrettifyEntry(record)

	c.Assert(string(bytes1), Equals, string(bytes2))
}

func (s *JsonPrettifierSuite) TestPrettifyEntry(c *C) {
	config.EnableLOC = true
	// The line number of below line is 34 which will be used as value of 'Line' field in record
	record := NewRecord(LOG_INFO, "Hello, world", map[string]string{"foo": "bar"})
	config.EnableLOC = false

	prettifier := NewJsonPrettifier(EXCLUDE_NONE)
	b, _ := prettifier.PrettifyEntry(record)

	// One example:
	// INFO 2012-09-27 16:53:40 json_prettifier_test.go:34:TestPrettifyEntry {"foo":"bar"} Hello, world
	c.Assert(string(b), Matches, `INFO .*son_prettifier_test.go:34:TestPrettifyEntry.*{"foo":"bar"}.*Hello, world`)
}

func (s *JsonPrettifierSuite) TestDecodeLogEntry(c *C) {
	config.EnableLOC = true
	// The line number of below line is 48 which will be used as value of 'Line' field in record
	record := NewRecord(LOG_INFO, "Hello, world", map[string]string{"foo": "bar"})
	config.EnableLOC = false
	record.Timestamp = 1348736601
	b, _ := NewJsonCodec().EncodeRecord(record)
	entry := string(b)

	prettifier := NewJsonPrettifier(EXCLUDE_NONE)
	record, err := prettifier.DecodeJsonLogEntry(entry)

	c.Assert(err, IsNil)
	c.Assert(record.Timestamp, Equals, int64(1348736601))
	c.Assert(record.Line, Equals, 48)
	c.Assert(record.Level, Equals, LOG_INFO)
	c.Assert(record.Method, Matches, ".*TestDecodeLogEntry$")
	c.Assert(record.Message, Equals, "Hello, world")
	c.Assert(record.File, Matches, ".*json_prettifier_test.go")
	c.Assert(record.Data["foo"], Equals, "bar")
}
