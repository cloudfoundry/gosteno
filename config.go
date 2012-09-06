package steno

type Config struct {
	Sinks []Sink
	Level *LogLevel
	Codec Codec
	Port  int

	EnableLOC bool

	User     string
	Password string
}
