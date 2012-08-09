package steno

import (
	"os"
	"testing"
)

const DEFAULT_LEVEL = "info"

func init() {
	jsonCodec := NewJsonCodec()
	config := NewConfig([]Sink{NewIOSink(os.Stdout)}, DEFAULT_LEVEL, jsonCodec)
	Init(config)
}

func TestLoggersNum(t *testing.T) {
	if len(loggers) != 0 {
		t.Error("The number of loggers at first should be ZERO.")
	}
}

func TestLoggerName(t *testing.T) {
	logger := NewLogger("test").(*BaseLogger)
	if logger.name != "test" {
		t.Error("It should return a logger with the name 'test'")
	}
}

func TestLoggerLevelActive(t *testing.T) {
	logger := NewLogger("foo").(*BaseLogger)
	defaultLevel := NewLogLevel(DEFAULT_LEVEL, 15)
	if !logger.active(defaultLevel) {
		t.Error("The default level " + DEFAULT_LEVEL + " should be activated")
	}

	higherLevels := []string{"warn", "error", "fatal"}
	for _, v := range higherLevels {
		level := NewLogLevel(v, 15)
		if !logger.active(level) {
			t.Error("The level '" + v + "' should be active for its priority is higher thant the default level : " + DEFAULT_LEVEL)
		}
	}
}

//FIXME:wait for the level modification api
func TestLevelModification(t *testing.T) {
}

func TestCreatingDupLogger(t *testing.T) {
	logger1 := NewLogger("foo")
	logger2 := NewLogger("foo")
	if logger1 != logger2 {
		t.Error("It should not create any new logger with an existing name")
	}
}
