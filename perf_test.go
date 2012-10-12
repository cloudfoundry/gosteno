package steno

import (
	"testing"
)

func BenchmarkNoSink(b *testing.B) {
	b.StopTimer()
	Init(&Config{})
	loggers = make(map[string]*BaseLogger)
	logger := NewLogger("nosink")
	b.StartTimer()

	performBenchmark(logger, b.N)
}

func BenchmarkDevNullSink(b *testing.B) {
	b.StopTimer()
	Init(&Config{
		Sinks: []Sink{NewFileSink("/dev/null")},
	})
	loggers = make(map[string]*BaseLogger)
	logger := NewLogger("dev_null_sink")
	b.StartTimer()

	performBenchmark(logger, b.N)
}

func BenchmarkDevNullSinkWithLOC(b *testing.B) {
	b.StopTimer()
	Init(&Config{
		Sinks:     []Sink{NewFileSink("/dev/null")},
		EnableLOC: true,
	})
	loggers = make(map[string]*BaseLogger)
	logger := NewLogger("dev_null_sink_with_loc")
	b.StartTimer()

	performBenchmark(logger, b.N)
}

func BenchmarkTaggedLoggerInDevNullSink(b *testing.B) {
	b.StopTimer()
	Init(&Config{
		Sinks: []Sink{NewFileSink("/dev/null")},
	})
	loggers = make(map[string]*BaseLogger)
	tags := map[string]string{
		"thread_id":    "1234",
		"program_name": "benchmark",
	}
	logger := NewTaggedLogger(NewLogger("dev_null_sink_tagged"), tags)
	b.StartTimer()

	performBenchmark(logger, b.N)
}

func performBenchmark(logger Logger, times int) {
	for i := 0; i < times; i++ {
		logger.Fatal("Hello, world.")
	}
}
