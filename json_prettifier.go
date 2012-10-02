package steno

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const (
	EXCLUDE_NONE = 0

	EXCLUDE_LEVEL = 1 << (iota - 1)
	EXCLUDE_TIMESTAMP
	EXCLUDE_FILE
	EXCLUDE_LINE
	EXCLUDE_METHOD
	EXCLUDE_DATA
	EXCLUDE_MESSAGE
)

type JsonPrettifier struct {
	entryTemplate *template.Template
}

func NewJsonPrettifier(flag int) *JsonPrettifier {
	fields := []string{
		"{{encodeLevel .Level}}",
		"{{encodeTimestamp .Timestamp}}",
		"{{encodeFile .File}}",
		"{{encodeLine .Line}}",
		"{{encodeMethod .Method}}",
		"{{encodeData .Data}}",
		"{{encodeMessage .Message}}",
	}

	for i, _ := range fields {
		// the shift count must be an unsigned integer
		if (flag & (1 << uint(i))) != 0 {
			fields[i] = ""
		}
	}

	prettifier := new(JsonPrettifier)
	format := strings.Join(fields, "")
	funcMap := template.FuncMap{
		"encodeTimestamp": encodeTimestamp,
		"encodeFile":      encodeFile,
		"encodeMethod":    encodeMethod,
		"encodeLine":      encodeLine,
		"encodeData":      encodeData,
		"encodeLevel":     encodeLevel,
		"encodeMessage":   encodeMessage,
	}
	prettifier.entryTemplate = template.Must(template.New("EntryTemplate").Funcs(funcMap).Parse(format))

	return prettifier
}

func (p *JsonPrettifier) DecodeLogEntry(logEntry string) (*Record, error) {
	var fieldsMap map[string]string
	if err := json.Unmarshal([]byte(logEntry), &fieldsMap); err != nil {
		return nil, err
	}

	record := new(Record)
	timestamp, err := time.Parse(TIME_FORMAT, fieldsMap["timestamp"])
	if err != nil {
		return nil, err
	}
	record.Timestamp = timestamp

	record.File = fieldsMap["file"]
	record.Method = fieldsMap["method"]
	record.Line, _ = strconv.Atoi(fieldsMap["line"])
	record.Level = LEVELS[fieldsMap["log_level"]]
	record.Message = fieldsMap["message"]

	var fields = map[string]bool{
		"timestamp": true,
		"file":      true,
		"method":    true,
		"line":      true,
		"log_level": true,
		"message":   true,
	}
	data := make(map[string]string)
	for k, v := range fieldsMap {
		if !fields[k] {
			data[k] = v
		}
	}
	record.Data = data

	return record, nil
}

func (p *JsonPrettifier) PrettifyEntry(record *Record) ([]byte, error) {
	buffer := bytes.NewBufferString("")
	err := p.entryTemplate.Execute(buffer, record)
	return buffer.Bytes(), err
}

func encodeLevel(level *LogLevel) string {
	return fmt.Sprintf("%s ", strings.ToUpper(level.String()))
}

func encodeTimestamp(t time.Time) string {
	return fmt.Sprintf("%s ", t.Format(TIME_FORMAT))
}

func encodeFile(file string) string {
	index := strings.LastIndex(file, "/")
	return fmt.Sprintf("%s:", file[index+1:])
}

func encodeLine(line int) string {
	return fmt.Sprintf("%s:", strconv.Itoa(line))
}

func encodeMethod(method string) string {
	index := strings.LastIndex(method, ".")
	return fmt.Sprintf("%s ", method[index+1:])
}

func encodeData(data map[string]string) (string, error) {
	bytes, err := json.Marshal(data)
	return fmt.Sprintf("%s ", string(bytes)), err
}

func encodeMessage(message string) string {
	return message
}
