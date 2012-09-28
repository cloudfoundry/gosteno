package steno

import (
	"encoding/json"
	"time"
)

const TIME_FORMAT = time.RFC1123

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
