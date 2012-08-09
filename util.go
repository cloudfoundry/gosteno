package steno

import "encoding/json"

func genErrorMsgInJson(err error) []byte {
	bytes, _ := json.Marshal(map[string]string{"error": err.Error()})
	return bytes
}
