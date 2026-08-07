package main

import (
	"dt2026/api"
	"dt2026/env"
	"log/slog"
	"net/http"
	"os"
)

var mustHaveEnvNames = []string{
	"DT_ENV",
}

func main() {
	slog.Info("Starting")

	_ = env.LoadEnv(".env")
	if err := env.CheckEnv(mustHaveEnvNames); err != nil {
		slog.Error("check env", "err", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", api.Health)

	server := http.Server{
		Addr:    ":29300",
		Handler: mux,
	}

	server.ListenAndServe()
}
