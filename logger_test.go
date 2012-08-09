package steno

import (
	"testing"
)

func setup() {

}

func teardown() {
  //The loggers should be clean up after every case done because it is a global variable
  //Let the gc clean the old loggers map
	loggers = make(map[string]Logger)
}

func TestLoggersNum(t *testing.T) {
  defer teardown()

	if len(loggers) != 0 {
		t.Error("The number of loggers at first should be ZERO.")
	}
}

func TestLoggerName(t *testing.T) {
  defer teardown()

	logger := NewLogger("test").(*BaseLogger)
	if logger.name != "test" {
		t.Error("It should return a logger with the name 'test'")
	}
}

func TestLoggerLevelActive(t *testing.T) {
  defer teardown()

	logger := NewLogger("foo").(*BaseLogger)
	defaultLevel := LOG_INFO
	if !logger.active(defaultLevel) {
		t.Error("The default level " + defaultLevel.name + " should be activated")
	}

	higherLevels := []*LogLevel{LOG_WARN, LOG_ERROR, LOG_ERROR}
	for _, level := range higherLevels {
		if !logger.active(level) {
			t.Error("The level '" + level.name + "' should be active for its priority is higher thant the default level : " + defaultLevel.name)
		}
	}
}

//FIXME:wait for the level modification api
func TestLevelModification(t *testing.T) {
  defer teardown()

}

func TestCreatingDupLogger(t *testing.T) {
  defer teardown()

	logger1 := NewLogger("foo")
	logger2 := NewLogger("foo")
	if logger1 != logger2 {
		t.Error("It should not create any new logger with an existing name")
	}
}
