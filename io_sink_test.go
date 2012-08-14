package steno

import (
	"bufio"
	"io"
	. "launchpad.net/gocheck"
	"os"
)

type IOSinkSuite struct {
}

var _ = Suite(&IOSinkSuite{})

func (s *IOSinkSuite) TestAddRecord(c *C) {
	pReader, pWriter := io.Pipe()
	sink := NewIOSink(nil)
	sink.writer = bufio.NewWriter(pWriter)
	sink.SetCodec(NewJsonCodec())

	go func(msg string) {
		record := NewRecord(LOG_INFO, msg, map[string]string{})
		sink.AddRecord(record)
		sink.Flush()
		pWriter.Close()
	}("Hello, \nworld")

	bufReader := bufio.NewReader(pReader)
	msg, err := bufReader.ReadString('\n')
	c.Assert(err, IsNil)
	c.Assert(msg, Matches, `{.*"Hello, \\nworld".*}\n`)
}

func (s *IOSinkSuite) TestMarshalJSON(c *C) {
	sink := NewIOSink(os.Stdout)
	sink.SetCodec(NewJsonCodec())
	msgBytes, _ := sink.MarshalJSON()
	c.Assert(string(msgBytes), Matches, `{.*"/dev/stdout".*}`)
}
