package records

import (
	"dt2026/httpx"
	"dt2026/lib"
	"dt2026/lib/optional"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
)

// JSONTime 序列化为 JSON 后总是 +08:00 时区
type JSONTime time.Time

// JSONSchemaTypes 给 jsonschema.For 用。JSONTime 是 time.Time 的定义类型，
// JSON 编成 RFC3339 字符串，但 For 只特殊处理 time.Time，否则会推断成空 object。
var JSONSchemaTypes = map[reflect.Type]*jsonschema.Schema{
	reflect.TypeFor[JSONTime](): {Type: "string"},
}

var tzSpace = regexp.MustCompile(`^(\S+) (\d{2}:\d{2})$`)

// ParseJSONTime 优先为 RFC3339Nano 格式，也支持本地日期
func ParseJSONTime(s string) (JSONTime, *httpx.Error) {
	// 时区的 X+08:00 在 URL query 中会被解析为 X 08:00
	s = tzSpace.ReplaceAllString(s, `$1+$2`)

	parsed, err := time.Parse(
		lib.RFC3339Nano, s,
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

// MarshalJSON 返回毫秒
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

	parsed, err := ParseJSONTime(string(data))

	if err != nil {
		return errors.New(err.Message)
	}

	*t = JSONTime(parsed)
	return nil
}
