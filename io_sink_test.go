package steno

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestNewIOSink(t *testing.T) {
	stdoutSink := NewIOSink(os.Stdout)
	if stdoutSink.writer == nil || stdoutSink.codec != nil {
		t.Error("It should return a steno.IO with the codec set to nil")
	}
}

func TestNewFileSink(t *testing.T) {
	logFile, _ := ioutil.TempFile("", "gosteno_test")
	defer os.Remove(logFile.Name())

	fileSink := NewIOSink(logFile)
	fileSink.SetCodec(NewJsonCodec())
	record := createTestRecord("Hello, world")
	fileSink.AddRecord(record)
	fileSink.Flush()

	reader := bufio.NewReader(logFile)
	logFile.Seek(0, 0)
	msg, _ := reader.ReadString('\n')
	if !(strings.HasPrefix(msg, "{") && strings.HasSuffix(msg, "}\n")) {
		t.Error("It should write a json string to the log file")
	}
}
