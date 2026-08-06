package api

import (
	"encoding/json"
	"net/http"
	"time"
)

var loc = time.FixedZone("CST", 8*3600)

func Health(w http.ResponseWriter, r *http.Request) {
	// 2026-08-06T22:39:47+08:00
	now := time.Now().In(loc).Format(time.RFC3339)
	err := json.NewEncoder(w).Encode(map[string]any{
		"ok":  true,
		"now": now,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ok":    false,
			"error": err.Error(),
		})
	}
}
