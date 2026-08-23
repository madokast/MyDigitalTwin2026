package probe

import (
	"dt2026/httpx"
	"net/http"
)

type BadJSONResponse struct {
	Ok     httpx.TrueType `json:"ok"`
	Status int         `json:"status"`
	Chan   chan string `json:"chan"` // This field is intentionally invalid for JSON marshalling
}

func (r BadJSONResponse) GetStatus() int {
	return r.Status
}

func BadJSON() BadJSONResponse {
	return BadJSONResponse{
		Ok:     httpx.True,
		Status: http.StatusOK,
		Chan:   make(chan string), // This will cause JSON marshalling to fail
	}
}
