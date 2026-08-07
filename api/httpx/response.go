package httpx

import (
	"encoding/json"
	"net/http"
)

func WriteJSONError(w http.ResponseWriter, err *Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status)

	json.NewEncoder(w).Encode(err)
}
