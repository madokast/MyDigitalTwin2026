package probe

import (
	"dt2026/lib"
	"net/http"
	"time"
)

const HealthDescription = "Probe health status of the MyDigitalTwin2026 server"

type HealthResponse struct {
	Ok     bool   `json:"ok" jsonschema:"请求相应结果"`
	Status int    `json:"status" jsonschema:"HTTP status code"`
	Now    string `json:"now" jsonschema:"当前服务器时间，采用 RFC3339Milli 格式字符串返回"`
}

func (r HealthResponse) GetStatus() int {
	return http.StatusOK
}

func Health() HealthResponse {
	// 2026-08-06T22:39:47+08:00
	now := time.Now().In(lib.UTC8).Format(lib.RFC3339Milli)
	return HealthResponse{
		Ok:     true,
		Status: http.StatusOK,
		Now:    now,
	}
}
