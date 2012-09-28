package steno

import "encoding/json"

func genErrorMsgInJson(err error) []byte {
	b, _ := json.Marshal(map[string]string{"error": err.Error()})
	return b
}
