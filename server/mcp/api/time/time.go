package time

import (
	"dt2026/api/probe"
	"dt2026/server/http"
)

const Name = "time"

const Description = "读取服务器当前本地时间（Asia/Shanghai，+08:00）。返回 RFC3339 毫秒时间、时区，以及中文日期、时刻和星期。需要当前时刻做上下文时用它，不要用它探活或读写记录。"

type Input struct {
}

type Output = probe.TimeResponse

func Time(_ *http.Server, _ Input) Output {
	return probe.Time()
}
