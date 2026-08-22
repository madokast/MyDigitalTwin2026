package probe

import (
	"dt2026/lib"
	"net/http"
	"time"
)

type HealthResponse struct {
	Ok     bool   `json:"ok"`
	Status int    `json:"status"`
	Now    string `json:"now"`
}

func (r HealthResponse) GetStatus() int {
	return http.StatusOK
}

func Health() HealthResponse {
	// 2026-08-06T22:39:47+08:00
	now := time.Now().In(lib.UTC8).Format(time.RFC3339)
	return HealthResponse{
		Ok:     true,
		Status: http.StatusOK,
		Now:    now,
	}
}
