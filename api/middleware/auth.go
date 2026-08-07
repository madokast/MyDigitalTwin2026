package middleware

import (
	"dt2026/api/envkeys"
	"dt2026/env"
	"dt2026/httpx"
	"net/http"
	"strings"
)

type handleFunc = func(w http.ResponseWriter, r *http.Request)

func auth(next handleFunc) handleFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := env.Get(envkeys.Token)
		if !ok {
			httpx.WriteJSON(w, httpx.NewInternalServerError("token not set"))
			return
		}

		auth := r.Header.Get("Authorization")

		if auth == "" {
			httpx.WriteJSON(w, httpx.NewUnauthorizedError("missing authorization header"))
			return
		}

		const prefix = "Bearer "

		if !strings.HasPrefix(auth, prefix) {
			httpx.WriteJSON(w, httpx.NewUnauthorizedError("invalid authorization header"+auth))
			return
		}

		provided := strings.TrimSpace(
			strings.TrimPrefix(auth, prefix),
		)

		if provided != token {
			httpx.WriteJSON(w, httpx.NewUnauthorizedError("invalid token "+provided))
			return
		}

		next(w, r)
	})
}
