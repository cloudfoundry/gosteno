package steno

import (
	"bufio"
	"os"
	"fmt"
)

type IO struct {
	writer *bufio.Writer
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
	// TODO: unified format
	msg := fmt.Sprintf("[%s] %s %s\n", record.timestamp, record.level.name, record.message)
	io.writer.WriteString(msg)
}


func (io *IO) Flush() {
	io.writer.Flush()
}
