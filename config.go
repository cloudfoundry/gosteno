package steno

type Config struct {
	sinks []Sink
	level *LogLevel
	codec Codec
	port  int
}

func NewConfig(sinks []Sink, level string, codec Codec, port int) *Config {
	var s = new(Config)

	s.level = lookupLevel(level)
	if s.level == nil {
		panic("Unknown level: " + level)
	}

	s.sinks = sinks
	s.codec = codec
	s.port = port

	for _, sink := range s.sinks {
		sink.SetCodec(codec)
	}

	return s
}
