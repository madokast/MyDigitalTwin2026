package health

import (
	"dt2026/api/probe"
	"dt2026/server/http"
)

const Name = "probe_health"

const Description = "探测服务器运行监控情况"

type Input struct {
}

type Output = probe.HealthResponse

func ProbeHealth(_ *http.Server, ph Input) Output {
	return probe.Health()
}
