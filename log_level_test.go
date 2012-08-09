package steno

import (
	"testing"
)

func TestNewLogLevel(t *testing.T) {
	level := NewLogLevel("foobar", 100)
	if (level == nil) || (level.name != "foobar") || (level.priority != 100) {
		t.Error("It should return a level with the name 'foobar' and priority set to 100")
	}
}

func TestLookupLevel(t *testing.T) {
	infoLevel := lookupLevel("info")
	if (infoLevel == nil) || (infoLevel.name != "info") || (infoLevel.priority != 15) {
		t.Error("It should return a level with the name 'info' and priority set to 15")
	}

	notExistLevel := lookupLevel("foobar")
	if notExistLevel != nil {
		t.Error("It should return a null level")
	}
}
