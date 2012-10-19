package steno

import (
	. "launchpad.net/gocheck"
	"strings"
)

type SyslogSinkSuite struct {
}

var _ = Suite(&SyslogSinkSuite{})

func (s *SyslogSinkSuite) TestTruncate(c *C) {
	msg := generateString(MaxMessageSize - 1)
	record := &Record{
		Message: msg,
	}
	truncate(record)
	c.Check(record.Message, Equals, msg)

	msg2 := generateString(MaxMessageSize + 1)
	record2 := &Record{
		Message: msg2,
	}
	truncate(record2)
	c.Check(strings.HasSuffix(record2.Message, "..."), Equals, true)
}

func generateString(length int) string {
	buffer := make([]byte, length)
	for i := 0; i < length; i++ {
		buffer[i] = byte('1')
	}
	return string(buffer)
}
