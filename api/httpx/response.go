package httpx

import (
	"encoding/json"
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
		w.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Sprintf(`{"ok":false,"status":500, "error":"failed to marshal response: %s"}`, err.Error())
		_, _ = w.Write([]byte(msg))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.GetStatus())
	_, err = w.Write(bs)
	if err != nil {
		slog.Error("failed to write response", "err", err)
	}
}
