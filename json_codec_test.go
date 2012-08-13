package steno

import (
	"reflect"
	"strings"
	"testing"
)

func createTestRecord(message string) *Record {
	data := make(map[string]string)
	return NewRecord(LOG_INFO, message, data)
}

func TestJsonCodec(t *testing.T) {
	record := createTestRecord("Hello world")

	codec := NewJsonCodec()
	msg := codec.EncodeRecord(record)
	if reflect.TypeOf(msg).String() != "string" {
		t.Error("The encoder should return a string")
	}
}

func TestEncodedResult(t *testing.T) {
	msg := "Hello, world"
	record := createTestRecord(msg)
	codec := NewJsonCodec()
	encodedRecord := codec.EncodeRecord(record)
	if !(strings.HasPrefix(encodedRecord, "{") && strings.HasSuffix(encodedRecord, "}\n")) {
		t.Error("It should return a json object string with a new line seperator")
	}
}

func TestEscapeNewLine(t *testing.T) {
	msg := "Newline\ntest"
	record := createTestRecord(msg)
	codec := NewJsonCodec()
	encodedRecord := codec.EncodeRecord(record)
	if !strings.Contains(encodedRecord, "ne\\nte") {
		t.Error("It should escape the new line")
	}
}

func TestEscapeCarriage(t *testing.T) {
	msg := "Newline\rtest"
	record := createTestRecord(msg)
	codec := NewJsonCodec()
	encodedRecord := codec.EncodeRecord(record)
	if !strings.Contains(encodedRecord, "ne\\rte") {
		t.Error("It should escape the carriage returns")
	}
}
