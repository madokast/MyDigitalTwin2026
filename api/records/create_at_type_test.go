package records

import (
	"strings"
	"testing"
	"time"
)

func TestMarshalJSON(t *testing.T) {
	createAt := JSONTime(time.Now())

	json, err := createAt.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(json), "+08:00") {
		t.Fatal("timezone not set", string(json))
	}
}

func TestUnmarshalJSON(t *testing.T) {
	json := `"2023-01-01T00:00:00+08:00"`

	var createAt JSONTime
	err := createAt.UnmarshalJSON([]byte(json))
	if err != nil {
		t.Fatal(err)
	}

	if createAt.Time().IsZero() {
		t.Fatal("time not set")
	}
}

func TestJsonMarshalJSONAndUnmarshalJSON(t *testing.T) {

	createAt := JSONTime(time.Now())

	json, err := createAt.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	var createAt2 JSONTime
	err = createAt2.UnmarshalJSON(json)
	if err != nil {
		t.Fatal(err)
	}
}
