package steno

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"
)

type TextCodec struct {
	entryTemplate *template.Template
}

func NewTextCodec(format string) Codec {
	textCodec := new(TextCodec)
	if strings.Contains(format, "{{.Data}}") {
		format = strings.Replace(format, "{{.Data}}", "{{encodeData .Data}}", -1)
	}
	funcMap := template.FuncMap{
		"encodeData": encodeData,
	}
	textCodec.entryTemplate = template.Must(template.New("EntryTemplate").Funcs(funcMap).Parse(format))
	return textCodec
}

func (t *TextCodec) EncodeRecord(record *Record) ([]byte, error) {
	buffer := bytes.NewBufferString("")
	err := t.entryTemplate.Execute(buffer, record)
	return buffer.Bytes(), err
}

func encodeData(data map[string]string) (string, error) {
	bytes, err := json.Marshal(data)
	return string(bytes), err
}
