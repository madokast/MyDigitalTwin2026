package middleware

import (
	"dt2026/api/envkeys"
	"dt2026/env"
	"net/http"
	"strings"
)

type handleFunc = func(w http.ResponseWriter, r *http.Request)

func Auth(next handleFunc) handleFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := env.Get(envkeys.Token)
		if !ok {
			http.Error(w, "token not set", http.StatusInternalServerError)
			return
		}

		auth := r.Header.Get("Authorization")

		const prefix = "Bearer "

		if !strings.HasPrefix(auth, prefix) {
			http.Error(w, "unauthorized", 401)
			return
		}

		provided := strings.TrimSpace(
			strings.TrimPrefix(auth, prefix),
		)

		if provided != token {
			http.Error(w, "unauthorized", 401)
			return
		}

		next(w, r)
	})
}
