package steno

import (
  "testing"
)

func TestNewRecord(t *testing.T) {
  message := "Hello, GOSTENO"
  data := make(map[string]string)
  record := NewRecord(LOG_INFO, message, data)

  if record == nil {
    t.Error("It should return a record")
  }
}
