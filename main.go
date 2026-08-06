package main

import (
	"dt2026/api"
	"log/slog"
	"net/http"
)

func main() {
	slog.Info("Starting")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", api.Health)

	server := http.Server{
		Addr:    ":29300",
		Handler: mux,
	}

	server.ListenAndServe()
}
