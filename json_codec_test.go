package steno

import (
	. "launchpad.net/gocheck"
	"strings"
)

func createTestRecord(message string) *Record {
	data := make(map[string]string)
	return NewRecord(LOG_INFO, message, data)
}

type JsonCodecSuite struct {
}

var _ = Suite(&JsonCodecSuite{})

func (s *JsonCodecSuite) SetUpTest(c *C) {

}

func (s *JsonCodecSuite) TearDownTest(c *C) {

}

func (s *JsonCodecSuite) TestJsonCodec(c *C) {
	record := createTestRecord("Hello world")

	codec := NewJsonCodec()
	msg, _ := codec.EncodeRecord(record)
	c.Assert(len(msg) > 0, Equals, true)
	c.Assert(msg, FitsTypeOf, []uint8{})
}

func (s *JsonCodecSuite) TestEncodedResult(c *C) {
	msg := "Hello, world"
	record := createTestRecord(msg)
	codec := NewJsonCodec()
	recordBytes, _ := codec.EncodeRecord(record)
	encodedRecord := string(recordBytes)
	c.Assert(encodedRecord, Matches, "{.*}")
}

func (s *JsonCodecSuite) TestEscapeNewLine(c *C) {
	msg := "Newline\ntest"
	record := createTestRecord(msg)
	codec := NewJsonCodec()
	recordBytes, _ := codec.EncodeRecord(record)
	encodedRecord := string(recordBytes)
	c.Assert(strings.Contains(encodedRecord, "Newline\\ntest"), Equals, true)
}

func (s *JsonCodecSuite) TestEscapeCarriage(c *C) {
	msg := "Newline\rtest"
	record := createTestRecord(msg)
	codec := NewJsonCodec()
	recordBytes, _ := codec.EncodeRecord(record)
	encodedRecord := string(recordBytes)
	c.Assert(strings.Contains(encodedRecord, "Newline\\rtest"), Equals, true)
}
