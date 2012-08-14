package steno

import (
	. "launchpad.net/gocheck"
)

func createTestRecord(message string) *Record {
	data := make(map[string]string)
	return NewRecord(LOG_INFO, message, data)
}

type JsonCodecSuite struct {
}

var _ = Suite(&JsonCodecSuite{})

func (s *JsonCodecSuite) TestJsonCodec(c *C) {
	record := createTestRecord("Hello world")
	codec := NewJsonCodec()
	msg, _ := codec.EncodeRecord(record)
	c.Assert(string(msg), Matches, `{.*"message":"Hello world".*}`)
}
