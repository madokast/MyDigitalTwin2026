package records

import (
	"dt2026/httpx"
	"fmt"
	"time"
)

// JSONTime 序列化为 JSON 后总是 +08:00 时区
type JSONTime time.Time

func ParseJSONTime(s string) (JSONTime, *httpx.Error) {
	parsed, err := time.Parse(
		time.RFC3339,
		s,
	)
	if err != nil {
		return JSONTime{}, httpx.NewBadRequestError(fmt.Sprintf(
			"invalid time value %s: %s",
			s,
			err.Error(),
		))
	}

	return JSONTime(parsed), nil
}

func (t *JSONTime) GoTime() time.Time {
	return time.Time(*t)
}

var loc = time.FixedZone("CST", 8*3600)

func (t JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(
		`"` + time.Time(t).In(loc).Format(time.RFC3339) + `"`,
	), nil
}

func (t *JSONTime) UnmarshalJSON(data []byte) error {
	if len(data) <= 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("invalid json-time value %s", data)
	}

	data = data[1 : len(data)-1]

	parsed, err := time.Parse(
		time.RFC3339,
		string(data),
	)
	if err != nil {
		return err
	}

	*t = JSONTime(parsed)
	return nil
}
