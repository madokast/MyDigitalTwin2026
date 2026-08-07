package httpx

import (
	"encoding/json/v2"
	"fmt"
	"log/slog"
	"net/http"
)

type Response interface {
	IsOk() bool
	GetStatus() int
}

func WriteJSON[R Response](w http.ResponseWriter, response R) {
	if response.GetStatus() == 0 {
		WriteJSON(w, NewInternalServerError("response status not set"))
		return
	}

	bs, err := json.Marshal(response)
	if err != nil {
		WriteJSON(w, NewInternalServerError(fmt.Sprintf("failed to marshal response: %s", err.Error())))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.GetStatus())
	_, err = w.Write(bs)
	if err != nil {
		slog.Error("failed to write response", "err", err)
	}
}
