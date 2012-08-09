package steno

import (
	"os"
	"testing"
)

func init() {
	jsonCodec := NewJsonCodec()
	config := NewConfig([]Sink{NewIOSink(os.Stdout)}, LOG_INFO.name, jsonCodec)
	Init(config)
}

func TestInit(t *testing.T) {
	defer teardown()

	if &config == nil {
		t.Error("The config should not be nil")
	}

	if len(loggers) != 0 {
		t.Error("The number of loggers at first should be ZERO.")
	}
}
