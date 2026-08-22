package health

import (
	"dt2026/api/probe"
	"dt2026/server/http"
)

const Name = "probe_health"

const Description = "探测服务是否可用，返回 ok、HTTP 状态和服务器当前时间（毫秒，+08:00）。不要用它查询历史或写入记录。"

type Input struct {
}

type Output = probe.HealthResponse

func ProbeHealth(_ *http.Server, ph Input) Output {
	return probe.Health()
}
