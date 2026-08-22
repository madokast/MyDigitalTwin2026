package probe

import (
	"dt2026/lib"
	"net/http"
	"time"
)

type HealthResponse struct {
	Ok     bool   `json:"ok" jsonschema:"是否成功"`
	Status int    `json:"status" jsonschema:"HTTP 状态码"`
	Now    string `json:"now" jsonschema:"服务器当前时间，毫秒精度，时区 +08:00"`
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
