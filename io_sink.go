package steno

import (
	"bufio"
	"os"
)

type IO struct {
	writer *bufio.Writer
	codec Codec
}

func NewIOSink(file *os.File) *IO {
	writer := bufio.NewWriter(file)

	io := new(IO)
	io.writer = writer

	return io
}

func NewFileSink(path string) *IO {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	return NewIOSink(file)
}

func (io *IO) AddRecord(record *Record) {
	msg := io.codec.EncodeRecord(record)
	io.writer.WriteString(msg)
}


func (io *IO) Flush() {
	io.writer.Flush()
}

func (io *IO) SetCodec(codec Codec) {
	io.codec = codec
}

func (io *IO) GetCodec() Codec {
	return io.codec
}
