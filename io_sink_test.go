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

func (s *IOSinkSuite) checkWriter(c *C, w *io.PipeWriter, msg string, ch chan int) {
	sink := NewIOSink(os.Stdout)
	sink.SetCodec(NewJsonCodec())
	sink.writer = bufio.NewWriter(w)
	record := createTestRecord(msg)
	sink.AddRecord(record)
	sink.Flush()
	w.Close()

	encodedMsg, _ := sink.codec.EncodeRecord(record)
	ch <- len(encodedMsg)
}

func (s *IOSinkSuite) TestAddRecord(c *C) {
	ch := make(chan int)
	pReader, pWriter := io.Pipe()

	go s.checkWriter(c, pWriter, "Hello, world", ch)

	var buf = make([]byte, 64)
	n := 0
	for {
		nn, err := pReader.Read(buf)
		if err == io.EOF {
			break
		}
		n += nn
	}
	c.Assert(<-ch, Equals, n-1)
}

func (s *IOSinkSuite) TestMarshalJSON(c *C) {
	sink := NewIOSink(os.Stdout)
	sink.SetCodec(NewJsonCodec())
	msgBytes, _ := sink.MarshalJSON()
	c.Assert(string(msgBytes), Equals, `{"type": "file", "file": "/dev/stdout"}`)
}
