package n26

import (
	"encoding/json"
	"testing"
	"time"
)

type TestData struct {
	Time  TimeStamp
	Other string
}

//2018-03-17 17:43:44.123 CET
const test = `{"Time": 1521308624123}`

var testTime = time.Date(2018, 03, 17, 16, 43, 44, 123*int(time.Millisecond), time.UTC)

func TestUnmarshal(t *testing.T) {
	a := TestData{}
	err := json.Unmarshal([]byte(test), &a)
	if err != nil {
		t.Error(err)
	}
	if !a.Time.Equal(testTime) {
		t.Error("Time does not match", a, testTime)
	}
}

func TestUnmarshalEmpty(t *testing.T) {
	a := TestData{}
	err := json.Unmarshal([]byte("{}"), &a)
	if err != nil {
		t.Error(err)
	}
	if a.Time != (TimeStamp{}) {
		t.Errorf("Time should not be set: %+v", a)
	}
}

var invalid = []string{`{"Time":1521308624123d}`, `{"Time":1521308s624123}`}

func TestUnmarshalError(t *testing.T) {
	for _, s := range invalid {
		a := TestData{}
		err := json.Unmarshal([]byte(s), &a)
		if a.Time != (TimeStamp{}) {
			t.Errorf("Time should not be set: %+v", a)
		}
		if err == nil {
			t.Error("Expected an error unmarshalling")
		}
	}
}
