package steno

import (
	"encoding/json"
	"strconv"
)

type JsonCodec struct {
}

func NewJsonCodec() Codec {
	return new(JsonCodec)
}

func (j *JsonCodec) EncodeRecord(record *Record) ([]byte, error) {
	hash := map[string]string{
		"Timestamp": record.Timestamp.Format("2006-01-02 15:04:05 -0700 MST"),
		"Message":   record.Message,
		"Log_level": record.Level.name,
	}

	if config.EnableLOC {
		hash["File"] = record.File
		hash["Method"] = record.Method
		hash["Line"] = strconv.Itoa(record.Line)
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
