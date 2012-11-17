package steno

type Codec interface {
	EncodeRecord(record *Record) ([]byte, error)
}
