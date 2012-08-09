package steno

type Config struct {
	sinks []Sink
	level string
	codec Codec
}

func NewConfig(sinks []Sink, level string, codec Codec) *Config {
	var s = new(Config)

	s.sinks = sinks
	s.level = level
	s.codec = codec

	for _, sink := range s.sinks {
		sink.SetCodec(codec)
	}

	return s;
}
