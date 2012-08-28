package steno

import (
	"syscall"
	"testing"
)

var LOG_FILE = "./test.log"

func BenchmarkNoSink(b *testing.B) {
	Init(&Config{})
	logger := NewLogger("nosink")

	performBenchmark(logger, b.N)
}

func BenchmarkDevNullSink(b *testing.B) {
	Init(&Config{
		Sinks: []Sink{NewFileSink("/dev/null")},
	})
	logger := NewLogger("dev_null_sink")

	performBenchmark(logger, b.N)
}

func BenchmarkDevNullSinkWithLOC(b *testing.B) {
	Init(&Config{
		Sinks:     []Sink{NewFileSink("/dev/null")},
		EnableLOC: true,
	})
	logger := NewLogger("dev_null_sink_with_loc")

	performBenchmark(logger, b.N)
}

func BenchmarkFileSink(b *testing.B) {
	Init(&Config{
		Sinks: []Sink{NewFileSink(LOG_FILE)},
	})
	logger := NewLogger("file_sink")

	performBenchmark(logger, b.N)

	b.StopTimer()
	syscall.Unlink(LOG_FILE)
}

func BenchmarkSyslogSink(b *testing.B) {
	Init(&Config{
		Sinks: []Sink{NewSyslogSink()},
	})
	logger := NewLogger("syslog")

	performBenchmark(logger, b.N)
}

func BenchmarkFileSinkWithLOC(b *testing.B) {
	Init(&Config{
		Sinks:     []Sink{NewFileSink(LOG_FILE)},
		EnableLOC: true,
	})
	logger := NewLogger("file_sink_with_loc")

	performBenchmark(logger, b.N)

	b.StopTimer()
	syscall.Unlink(LOG_FILE)
}

func BenchmarkTaggedLoggerInFileSink(b *testing.B) {
	Init(&Config{
		Sinks: []Sink{NewFileSink(LOG_FILE)},
	})

	tags := map[string]string{
		"thread_id":    "1234",
		"program_name": "benchmark",
	}
	logger := NewTaggedLogger(NewLogger("file_sink"), tags)

	performBenchmark(logger, b.N)

	b.StopTimer()
	syscall.Unlink(LOG_FILE)
}

func BenchmarkTaggedLoggerInSyslogSink(b *testing.B) {
	Init(&Config{
		Sinks: []Sink{NewSyslogSink()},
	})

	tags := map[string]string{
		"thread_id":    "1234",
		"program_name": "benchmark",
	}
	logger := NewTaggedLogger(NewLogger("syslog"), tags)

	performBenchmark(logger, b.N)
}

func performBenchmark(logger Logger, times int) {
	for i := 0; i < times; i++ {
		logger.Fatal("Hello, world.")
	}
}
