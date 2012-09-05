package steno

import (
	"fmt"
	. "launchpad.net/gocheck"
	"time"
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
	// The line number of below line is 36 which will be used as value of 'Line' field in record
	record := NewRecord(LOG_INFO, "Hello, world", map[string]string{"foo": "bar"})
	config.EnableLOC = false

	prettifier := NewJsonPrettifier(EXCLUDE_NONE)
	bytes, _ := prettifier.PrettifyEntry(record)
	// one instance :
	// Error: INFO 2012/09/10 14:17:12 json_prettifier_test.go:36:TestNewJsonPrettifier {"foo":"bar"} Hello, world
	c.Assert(string(bytes), Matches, `INFO .*son_prettifier_test.go:36:TestPrettifyEntry.*{"foo":"bar"}.*Hello, world.*`)
}

func (s *JsonPrettifierSuite) TestEncodeTimestamp(c *C) {
	t, _ := time.Parse("2006-01-02 15:04:05", "2012-09-10 12:00:00")
	str := fmt.Sprintf("%d/0%d/%d %d:0%d:0%d ", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	c.Assert(encodeTimestamp(t), Equals, str)
}

func (s *JsonPrettifierSuite) TestDecodeLogEntry(c *C) {
	entry := `{"File":"/mnt/hgfs/vmshared/workspace-microcloud/go/tutorial/gosteno/simple_steno.go",
  "Jeff":"Li","Line":"12","Log_level":"fatal","Message":"Fatal","Method":"main.keepWritingLogs",
  "Timestamp":"2012-09-10 17:07:47 +0800 CST","bar":"foo","foo":"bar"}`

	prettifier := NewJsonPrettifier(EXCLUDE_NONE)
	record, err := prettifier.DecodeLogEntry(entry)
	c.Assert(err, IsNil)
	c.Assert(record.Line, Equals, 12)
	c.Assert(record.Level, Equals, LOG_FATAL)
	c.Assert(record.Method, Equals, "main.keepWritingLogs")
	c.Assert(record.Message, Equals, "Fatal")
	c.Assert(record.File, Equals, "/mnt/hgfs/vmshared/workspace-microcloud/go/tutorial/gosteno/simple_steno.go")
	c.Assert(record.Data["foo"], Equals, "bar")
	c.Assert(record.Data["bar"], Equals, "foo")
}
