package steno

type Codec interface {
	EncodeRecord(record *Record) string
}
