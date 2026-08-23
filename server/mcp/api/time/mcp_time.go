package time

import (
	"dt2026/api/probe"
	"dt2026/httpx"
	"dt2026/server/http"
)

const Name = "time"

const Description = "Read the server's current local time (Asia/Shanghai, +08:00). Returns an RFC 3339 timestamp with milliseconds, the time zone, and the date, clock time, and weekday in Chinese. Use this when you need the current time as context; do not use it as a health check or to read or write records."

type Input struct {
}

type Output = probe.TimeResponse

func Time(_ *http.Server, _ Input) (Output, *httpx.Error) {
	return probe.Time(), nil
}
