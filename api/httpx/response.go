package httpx

import (
	"encoding/json"
	"net/http"
)

type Response interface {
	IsOk() bool
	GetStatus() int
}

func WriteJSON[R Response](w http.ResponseWriter, response R) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.GetStatus())

	_ = json.NewEncoder(w).Encode(response)
}
