package steno

import (
	. "launchpad.net/gocheck"
	"testing"
)

type TextCodecSuite struct {
}

var _ = Suite(&TextCodecSuite{})

func Test(t *testing.T) {
	TestingT(t)
}

func (s *TextCodecSuite) TestInvalidFormat(c *C) {
	c.Assert(func() { NewTextCodec("{{.Message}") }, PanicMatches, `.*unexpected "}".*`)
}

func (s *TextCodecSuite) TestTextCodecEncodeRecord(c *C) {
	codec := NewTextCodec("foo:{{.Message}}:bar")
	record := NewRecord(LOG_INFO, "Hello, world", map[string]string{})
	bytes, err := codec.EncodeRecord(record)
	c.Assert(err, IsNil)
	c.Assert(string(bytes), Equals, "foo:Hello, world:bar")
}

func (s *TextCodecSuite) TestNotStringFields(c *C) {
	codec := NewTextCodec("{{.Level}} : {{.Message}}")
	record := NewRecord(LOG_FATAL, "whatever", map[string]string{})
	bytes, err := codec.EncodeRecord(record)
	c.Assert(err, IsNil)
	c.Assert(string(bytes), Equals, "fatal : whatever")
}

func (s *TextCodecSuite) TestWithLOCDisable(c *C) {
	codec := NewTextCodec("{{.File}}:{{.Line}}:{{.Method}}")
	record := NewRecord(LOG_FATAL, "whatever", map[string]string{})
	bytes, err := codec.EncodeRecord(record)
	c.Assert(err, IsNil)
	c.Assert(string(bytes), Equals, ""+":0:"+"")
}

func (s *TextCodecSuite) TestWithData(c *C) {
	codec := NewTextCodec("{{.Message}}:{{.Data}}")
	record := NewRecord(LOG_FATAL, "Hello,world", map[string]string{"foo": "bar", "meta": "{{"})
	bytes, err := codec.EncodeRecord(record)
	c.Assert(err, IsNil)
	c.Assert(string(bytes), Equals, `Hello,world:{"foo":"bar","meta":"{{"}`)
}
