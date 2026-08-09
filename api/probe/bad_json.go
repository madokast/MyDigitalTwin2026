package probe

import (
	"net/http"
)

type BadJSONResponse struct {
	Ok     bool        `json:"ok"`
	Status int         `json:"status"`
	Chan   chan string `json:"chan"` // This field is intentionally invalid for JSON marshalling
}

func (r BadJSONResponse) GetStatus() int {
	return r.Status
}

func BadJSON() BadJSONResponse {
	return BadJSONResponse{
		Ok:     true,
		Status: http.StatusOK,
		Chan:   make(chan string), // This will cause JSON marshalling to fail
	}
}
