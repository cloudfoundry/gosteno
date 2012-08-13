package steno

import "encoding/json"

type Sink interface {
	json.Marshaler

	AddRecord(record *Record)
	Flush()

	SetCodec(codec Codec)
	GetCodec() Codec
}
