package records

import (
	"dt2026/httpx"
	"dt2026/lib"
	"dt2026/lib/optional"
	"fmt"
	"time"
)

// JSONTime 序列化为 JSON 后总是 +08:00 时区
type JSONTime time.Time

func ParseJSONTime(s string) (JSONTime, *httpx.Error) {
	parsed, err := time.Parse(
		lib.RFC3339Milli, s,
	)
	if err != nil {
		parsed, err = time.ParseInLocation(
			time.DateOnly, s, lib.UTC8,
		)

		if err != nil {
			return JSONTime{}, httpx.NewBadRequestError(fmt.Sprintf(
				"time value expected in RFC3339 or DateOnly, but got '%s'",
				s,
			))
		}
	}

	return JSONTime(parsed), nil
}

func ParseOptionalJSONTime(opt optional.Optional[string]) (optional.Optional[JSONTime], *httpx.Error) {
	s, ok := opt.Get()
	if !ok {
		return optional.None[JSONTime](), nil
	}

	parsed, err := ParseJSONTime(s)
	if err != nil {
		return optional.None[JSONTime](), err
	}

	return optional.Some(parsed), nil
}

func (t *JSONTime) GoTime() time.Time {
	return time.Time(*t)
}

func (t JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(
		`"` + time.Time(t).In(lib.UTC8).Format(lib.RFC3339Milli) + `"`,
	), nil
}

func (t *JSONTime) UnmarshalJSON(data []byte) error {
	if len(data) <= 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("invalid json-time value %s", data)
	}

	data = data[1 : len(data)-1]

	parsed, err := time.Parse(
		lib.RFC3339Milli,
		string(data),
	)
	if err != nil {
		return err
	}

	*t = JSONTime(parsed)
	return nil
}
