package steno

import (
	"testing"
	"fmt"
	"time"
	"syscall"
)

var LOG_FILE = "./test.log"

var testC = &Config{
	Sinks: []Sink{ nil },
	Level: LOG_INFO,
	Codec: JSON_CODEC,
}

func TestFileSink(t *testing.T) {
	testC.Sinks[0] = NewFileSink(LOG_FILE)
	Init(testC)
	logger := NewLogger("test")

	fmt.Println("# testing file sink:")

	performBenchmark(logger, 10000)

	syscall.Unlink(LOG_FILE)
}

func TestSyslogSink(t *testing.T) {
	testC.Sinks[0] = NewSyslogSink()
	Init(testC)
	logger := NewLogger("test")

	fmt.Println("# testing syslog sink:")

	performBenchmark(logger, 10000)
}

func TestFileSinkWithLOC(t *testing.T) {
	testC.Sinks[0] = NewFileSink(LOG_FILE)
	testC.EnableLOC = true
	Init(testC)
	logger := NewLogger("test")

	fmt.Println("# testing file sink with loc:")

	performBenchmark(logger, 10000)

	syscall.Unlink(LOG_FILE)
}

func TestTaggedLoggerInFileSink(t *testing.T) {
	testC.Sinks[0] = NewFileSink(LOG_FILE)
	Init(testC)

	tags := map[string]string {
		"thread_id": "1234",
		"program_name": "benchmark",
	}
	logger := NewTaggedLogger(NewLogger("test"), tags)

	fmt.Println("# testing tagged logger with file sink:")

	performBenchmark(logger, 10000)

	syscall.Unlink(LOG_FILE)
}

func TestTaggedLoggerInSyslogSink(t *testing.T) {
	testC.Sinks[0] = NewSyslogSink()
	Init(testC)

	tags := map[string]string {
		"thread_id": "1234",
		"program_name": "benchmark",
	}
	logger := NewTaggedLogger(NewLogger("test"), tags)

	fmt.Println("# testing tagged logger with syslog sink:")

	performBenchmark(logger, 10000)

	syscall.Unlink(LOG_FILE)
}

func performBenchmark(logger Logger, times int) {
	start := time.Now()
	fmt.Printf("writing %d logs ...\n", times)
	for i:= 1; i <= times; i++ {
		logger.Fatal("Hello, world.")
		if i % 500 == 0 {
			fmt.Print("+")
		}
	}
	spent := time.Since(start)
	rate := time.Duration(times) * time.Second / spent
	fmt.Println()

	fmt.Printf("time spent: %s, logs per second: %d\n\n", spent, rate)
}
