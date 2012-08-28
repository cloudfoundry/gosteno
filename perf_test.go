package steno

import (
	"testing"
)

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

func BenchmarkTaggedLoggerInDevNullSink(b *testing.B) {
	Init(&Config{
		Sinks: []Sink{NewFileSink("/dev/null")},
	})

	tags := map[string]string{
		"thread_id":    "1234",
		"program_name": "benchmark",
	}
	logger := NewTaggedLogger(NewLogger("dev_null_sink_tagged"), tags)

	performBenchmark(logger, b.N)
}

func performBenchmark(logger Logger, times int) {
	for i := 0; i < times; i++ {
		logger.Fatal("Hello, world.")
	}
}
