package probe

import (
	"dt2026/httpx"
	"dt2026/lib"
	"net/http"
	"time"
)

type HealthResponse struct {
	Ok     httpx.TrueType `json:"ok" jsonschema:"Whether the request succeeded"`
	Status int    `json:"status" jsonschema:"HTTP status code"`
	Now    string `json:"now" jsonschema:"Current server time in RFC 3339 with millisecond precision and a +08:00 offset"`
}

func (r HealthResponse) GetStatus() int {
	return http.StatusOK
}

func Health() HealthResponse {
	// 2026-08-06T22:39:47+08:00
	now := time.Now().In(lib.UTC8).Format(lib.RFC3339Milli)
	return HealthResponse{
		Ok:     httpx.True,
		Status: http.StatusOK,
		Now:    now,
	}
}
