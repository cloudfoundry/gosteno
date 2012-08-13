package steno

import (
	"bufio"
	"fmt"
	"os"
)

type IOSink struct {
	writer *bufio.Writer
	codec  Codec
	file   *os.File
}

func NewIOSink(file *os.File) *IOSink {
	writer := bufio.NewWriter(file)

	ioSink := new(IOSink)
	ioSink.writer = writer
	ioSink.file = file

	return ioSink
}

func NewFileSink(path string) *IOSink {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	return NewIOSink(file)
}

func (ioSink *IOSink) AddRecord(record *Record) {
	bytes := ioSink.codec.EncodeRecord(record)
	ioSink.writer.Write(bytes)

	// Need to append a newline for IO sink
	ioSink.writer.WriteString("\n")
}

func (ioSink *IOSink) Flush() {
	ioSink.writer.Flush()
}

func (ioSink *IOSink) SetCodec(codec Codec) {
	ioSink.codec = codec
}

func (ioSink *IOSink) GetCodec() Codec {
	return ioSink.codec
}

func (ioSink *IOSink) MarshalJSON() ([]byte, error) {
	msg := fmt.Sprintf("{\"type\": \"file\", \"file\": \"%s\"}", ioSink.file.Name())
	return []byte(msg), nil
}
