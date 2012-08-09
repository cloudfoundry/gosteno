package steno

type Config struct {
	sinks []Sink
	level *LogLevel
	codec Codec
}

func NewConfig(sinks []Sink, level string, codec Codec) *Config {
	var s = new(Config)

	s.level = lookupLevel(level)
	if s.level == nil {
		panic("Unknown level: " + level)
	}

	s.sinks = sinks
	s.codec = codec

	for _, sink := range s.sinks {
		sink.SetCodec(codec)
	}

	return s;
}
