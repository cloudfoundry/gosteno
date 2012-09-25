package steno

import (
	"encoding/json"
	"strconv"
	"time"
)

const TimeFormat = time.RFC1123

type JsonCodec struct {
}

func NewJsonCodec() Codec {
	return new(JsonCodec)
}

func (j *JsonCodec) EncodeRecord(record *Record) ([]byte, error) {
	hash := map[string]string{
		"timestamp": record.Timestamp.Format(TimeFormat),
		"message":   record.Message,
		"log_level": record.Level.name,
	}

	if config.EnableLOC {
		hash["file"] = record.File
		hash["method"] = record.Method
		hash["line"] = strconv.Itoa(record.Line)
	}

	if record.Data != nil {
		// Notice: it is possible data overwrite other record
		for k, v := range record.Data {
			hash[k] = v
		}
	}

	bytes, err := json.Marshal(hash)
	if err != nil {
		bytes = genErrorMsgInJson(err)
	}
	return bytes, err
}
