package main

import (
	"dt2026/api/envkeys"
	"dt2026/api/middleware"
	"dt2026/api/probe"
	"dt2026/env"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	slog.Info("Starting")

	_ = env.LoadEnv(".env")
	if err := env.CheckEnv(envkeys.All); err != nil {
		slog.Error("check env", "err", err)
		os.Exit(1)
	}
	slog.Info("Env loaded")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/probe/health", middleware.Auth(probe.Health))
	mux.HandleFunc("GET /api/probe/bad-json", middleware.Auth(probe.BadJSON))

	server := http.Server{
		Addr:    ":" + env.MustGet(envkeys.ServerPort),
		Handler: mux,
	}

	slog.Info("Server starting", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}
