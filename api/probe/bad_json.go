package probe

import (
	"dt2026/httpx"
	"net/http"
)

type BadJSONResponse struct {
	Ok     bool        `json:"ok"`
	Status int         `json:"status"`
	Chan   chan string `json:"chan"` // This field is intentionally invalid for JSON marshalling
}

func (r BadJSONResponse) IsOk() bool {
	return r.Ok
}

func (r BadJSONResponse) GetStatus() int {
	return r.Status
}

func BadJSON(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, BadJSONResponse{
		Ok:     true,
		Status: http.StatusOK,
		Chan:   make(chan string), // This will cause JSON marshalling to fail
	})
}
