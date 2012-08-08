package steno

type Config struct {
	sinks []Sink
	level string
}

func NewConfig(sinks []Sink, level string) *Config {
	var s = new(Config)

	s.sinks = sinks
	s.level = level

	return s;
}
