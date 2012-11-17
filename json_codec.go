package steno

import (
	"encoding/json"
)

type JsonCodec struct {
}

func NewJsonCodec() Codec {
	return new(JsonCodec)
}

func (j *JsonCodec) EncodeRecord(record *Record) ([]byte, error) {
	b, err := json.Marshal(record)
	if err != nil {
		b = genErrorMsgInJson(err)
	}

	return b, err
}
