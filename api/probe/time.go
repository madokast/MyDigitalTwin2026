package probe

import (
	"dt2026/lib"
	"fmt"
	"net/http"
	"time"
)

type TimeResponse struct {
	Ok       bool   `json:"ok" jsonschema:"Whether the request succeeded"`
	Status   int    `json:"status" jsonschema:"HTTP status code"`
	Datetime string `json:"datetime" jsonschema:"Current server time in RFC 3339 with millisecond precision and a +08:00 offset"`
	Timezone string `json:"timezone" jsonschema:"IANA time zone name; always Asia/Shanghai"`
	Local    Local  `json:"local" jsonschema:"Human-readable local date and time"`
}

type Local struct {
	Date    string `json:"date" jsonschema:"Date in Chinese, e.g. 2026年8月11日"`
	Time    string `json:"time" jsonschema:"Time of day in Chinese, e.g. 10点30分00秒"`
	Weekday string `json:"weekday" jsonschema:"Weekday in Chinese, e.g. 星期二"`
}

func (t TimeResponse) GetStatus() int {
	return t.Status
}

var chineseWeekdays = map[time.Weekday]string{
	time.Monday:    "星期一",
	time.Tuesday:   "星期二",
	time.Wednesday: "星期三",
	time.Thursday:  "星期四",
	time.Friday:    "星期五",
	time.Saturday:  "星期六",
	time.Sunday:    "星期日",
}

func Time() TimeResponse {
	now := time.Now().In(lib.UTC8)

	return TimeResponse{
		Ok:       true,
		Status:   http.StatusOK,
		Datetime: now.Format(lib.RFC3339Milli),
		Timezone: "Asia/Shanghai",
		Local: Local{
			Date: fmt.Sprintf(
				"%d年%d月%d日",
				now.Year(),
				now.Month(),
				now.Day(),
			),
			Time: fmt.Sprintf(
				"%d点%02d分%02d秒",
				now.Hour(),
				now.Minute(),
				now.Second(),
			),
			Weekday: chineseWeekdays[now.Weekday()],
		},
	}
}
