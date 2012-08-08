package steno

type Sink interface {
	AddRecord(record *Record)
	Flush()
}
