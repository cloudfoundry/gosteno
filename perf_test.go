package steno

import (
	"testing"
)

func BenchmarkNoSink(b *testing.B) {
	Init(&Config{})
	loggers = make(map[string]*BaseLogger)
	logger := NewLogger("nosink")

	performBenchmark(logger, b)
}

func BenchmarkDevNullSink(b *testing.B) {
	Init(&Config{
		Sinks: []Sink{NewFileSink("/dev/null")},
	})
	loggers = make(map[string]*BaseLogger)
	logger := NewLogger("dev_null_sink")

	performBenchmark(logger, b)
}

func BenchmarkDevNullSinkWithLOC(b *testing.B) {
	Init(&Config{
		Sinks:     []Sink{NewFileSink("/dev/null")},
		EnableLOC: true,
	})
	loggers = make(map[string]*BaseLogger)
	logger := NewLogger("dev_null_sink_with_loc")

	performBenchmark(logger, b)
}

func BenchmarkTaggedLoggerInDevNullSink(b *testing.B) {
	Init(&Config{
		Sinks: []Sink{NewFileSink("/dev/null")},
	})
	loggers = make(map[string]*BaseLogger)
	tags := map[string]string{
		"thread_id":    "1234",
		"program_name": "benchmark",
	}
	logger := NewTaggedLogger(NewLogger("dev_null_sink_tagged"), tags)

	performBenchmark(logger, b)
}

func performBenchmark(logger Logger, b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Fatal("Hello, world.")
	}
}
