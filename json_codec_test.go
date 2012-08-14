package steno

import (
	. "launchpad.net/gocheck"
)

type JsonCodecSuite struct {
}

var _ = Suite(&JsonCodecSuite{})

func (s *JsonCodecSuite) TestJsonCodec(c *C) {
	record := NewRecord(LOG_INFO, "Hello world", map[string]string{})
	codec := NewJsonCodec()
	msg, _ := codec.EncodeRecord(record)
	c.Assert(string(msg), Matches, `{.*"Message":"Hello world".*}`)
}
