package steno

import (
	"encoding/json"
	"fmt"
)

type JsonCodec struct {
}

func NewJsonCodec() Codec {
	return new(JsonCodec)
}

func (j *JsonCodec) EncodeRecord(record *Record) []byte {
	hash := map[string]string{
		"timestamp": record.timestamp.String(),
		"message":   record.message,
		"log_level": record.level.name,
	}

	if record.data != nil {
		// Notice: it is possible data overwrite other record
		for k, v := range record.data {
			hash[k] = v
		}
	}

	bytes, err := json.Marshal(hash)
	if err != nil {
		message := fmt.Sprintf("Error: Encoding JsonCodec, record: (%s)", err)
		bytes = []byte(message)
	}
	return bytes
}
