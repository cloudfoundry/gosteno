package steno

import (
	. "launchpad.net/gocheck"
)

type JsonPrettifierSuite struct {
}

var _ = Suite(&JsonPrettifierSuite{})

func (s *JsonPrettifierSuite) TestConst(c *C) {
	c.Assert(EXCLUDE_NONE, Equals, 0)
	c.Assert(EXCLUDE_LEVEL, Equals, 1<<0)
	c.Assert(EXCLUDE_TIMESTAMP, Equals, 1<<1)
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
	bytes, _ := prettifier.PrettifyEntry(record)

	// One example:
	// INFO Wed, 19 Sep 2012 10:51:57 CST json_prettifier_test.go:34:TestPrettifyEntry {"foo":"bar"} Hello, world
	c.Assert(string(bytes), Matches, `INFO .*son_prettifier_test.go:34:TestPrettifyEntry.*{"foo":"bar"}.*Hello, world`)
}

func (s *JsonPrettifierSuite) TestExclude(c *C) {
	config.EnableLOC = true
	record := NewRecord(LOG_INFO, "Hello, world", map[string]string{"foo": "bar"})
	config.EnableLOC = false

	prettifier := NewJsonPrettifier(EXCLUDE_DATA | EXCLUDE_LINE)
	bytes, _ := prettifier.PrettifyEntry(record)

	// One example:
	// INFO Wed, 19 Sep 2012 10:51:57 CST json_prettifier_test.go:TestExclude Hello, world
	c.Assert(string(bytes), Matches, `INFO .*son_prettifier_test.go:TestExclude Hello, world`)
}

func (s *JsonPrettifierSuite) TestDecodeLogEntry(c *C) {
	entry := `{"file":"/tmp/gopath/src/gosteno/json_prettifier_test.go","foo":"bar","line":"57",
  "log_level":"info","message":"Hello, world","method":"gosteno.(*JsonPrettifierSuite).TestDecodeLogEntry",
  "timestamp":"Wed, 19 Sep 2012 00:19:29 CST"}`

	prettifier := NewJsonPrettifier(EXCLUDE_NONE)
	record, err := prettifier.DecodeLogEntry(entry)

	c.Assert(err, IsNil)
	c.Assert(record.Line, Equals, 57)
	c.Assert(record.Level, Equals, LOG_INFO)
	c.Assert(record.Method, Equals, "gosteno.(*JsonPrettifierSuite).TestDecodeLogEntry")
	c.Assert(record.Message, Equals, "Hello, world")
	c.Assert(record.File, Equals, "/tmp/gopath/src/gosteno/json_prettifier_test.go")
	c.Assert(record.Data["foo"], Equals, "bar")
}
