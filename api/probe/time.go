package probe

import (
	"dt2026/lib"
	"fmt"
	"net/http"
	"time"
)

type TimeResponse struct {
	Ok       bool   `json:"ok"`
	Status   int    `json:"status"`
	Datetime string `json:"datetime"`
	Timezone string `json:"timezone"`
	Local    Local  `json:"local"`
}

type Local struct {
	Date    string `json:"date"`
	Time    string `json:"time"`
	Weekday string `json:"weekday"`
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
		Datetime: now.Format(time.RFC3339),
		Timezone: "Asia/Shanghai",
		Local: Local{
			Date: fmt.Sprintf(
				"%d年%d月%d日",
				now.Year(),
				now.Month(),
				now.Day(),
			),
			Time: fmt.Sprintf(
				"%d点%d分%d秒",
				now.Hour(),
				now.Minute(),
				now.Second(),
			),
			Weekday: chineseWeekdays[now.Weekday()],
		},
	}
}
