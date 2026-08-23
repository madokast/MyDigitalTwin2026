package probe

import (
	"dt2026/lib"
	"fmt"
	"net/http"
	"time"
)

type TimeResponse struct {
	Ok       bool   `json:"ok" jsonschema:"是否成功"`
	Status   int    `json:"status" jsonschema:"HTTP 状态码"`
	Datetime string `json:"datetime" jsonschema:"服务器当前时间，毫秒精度，RFC3339，时区 +08:00"`
	Timezone string `json:"timezone" jsonschema:"IANA 时区名，固定 Asia/Shanghai"`
	Local    Local  `json:"local" jsonschema:"面向人读的本地时间"`
}

type Local struct {
	Date    string `json:"date" jsonschema:"中文日期，如 2026年8月11日"`
	Time    string `json:"time" jsonschema:"中文时刻，如 10点30分00秒"`
	Weekday string `json:"weekday" jsonschema:"中文星期，如 星期二"`
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
