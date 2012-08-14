package steno

var JSON_CODEC = NewJsonCodec()

type Codec interface {
	EncodeRecord(record *Record) ([]byte, error)
}
