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

	EXCLUDE_TIMESTAMP = 1 << (iota - 1)
	EXCLUDE_FILE
	EXCLUDE_METHOD
	EXCLUDE_LINE
	EXCLUDE_DATA
	EXCLUDE_LEVEL
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

	for i := 0; i < len(fields); i++ {
		// an unsigned integer is required in shift operation
		if (flag & (1 << uint32(i))) != 0 {
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
	record := new(Record)
	if err := json.Unmarshal([]byte(logEntry), &fieldsMap); err != nil {
		return record, err
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05 -0700 MST", fieldsMap["Timestamp"])
	if err != nil {
		return record, err
	}
	record.Timestamp = timestamp
	record.File = fieldsMap["File"]
	record.Method = fieldsMap["Method"]
	record.Line, _ = strconv.Atoi(fieldsMap["Line"])
	record.Level = LEVELS[fieldsMap["Log_level"]]
	record.Message = fieldsMap["Message"]

	var fields = map[string]bool{
		"Timestamp": true,
		"File":      true,
		"Method":    true,
		"Line":      true,
		"Log_level": true,
		"Message":   true,
	}
	data := make(map[string]string)
	for k, v := range fieldsMap {
		if !fields[k] {
			data[k] = v
		}
	}
	record.Data = data

	return record, err
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
	return fmt.Sprintf("%s ", t.Format("2006/01/02 15:04:05"))
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
