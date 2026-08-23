package health

import (
	"dt2026/api/probe"
	"dt2026/httpx"
	"dt2026/server/http"
)

const Name = "probe_health"

const Description = "Probe whether the service is up. Returns ok, the HTTP status, and the current server time (milliseconds, +08:00). Do not use this to query history or write records."

type Input struct {
}

type Output = probe.HealthResponse

func ProbeHealth(_ *http.Server, _ Input) (Output, *httpx.Error) {
	return probe.Health(), nil
}
