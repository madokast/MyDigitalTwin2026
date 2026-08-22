package main

import (
	"dt2026/api/envkeys"
	"dt2026/api/middleware"
	"dt2026/env"
	http_server "dt2026/server/http"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	slog.Info("Starting")

	if err := env.LoadEnv(".env"); err != nil {
		slog.Error("failed to load env", "err", err)
		os.Exit(1)
	}
	server := http_server.NewServer()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/time", middleware.Auth(server.Time))
	mux.HandleFunc("GET /api/probe/health", middleware.Auth(server.ProbeHealth))
	mux.HandleFunc("GET /api/probe/bad-json", middleware.Auth(server.ProbeBadJson))
	mux.HandleFunc("GET /api/probe/postgresql", middleware.Auth(server.ProbePostgreSQL))
	mux.HandleFunc("GET /api/probe/qqbot", middleware.Auth(server.ProbeQQBot))
	mux.HandleFunc("POST /api/records", middleware.Auth(server.RecordsPost))
	mux.HandleFunc("GET /api/records", middleware.Auth(server.RecordsQuery))
	mux.HandleFunc("GET /api/records/export", middleware.Auth(server.RecordsExport))
	mux.HandleFunc("GET /api/records/{record_id}", middleware.Auth(server.RecordGet))
	mux.HandleFunc("GET /api/records/{record_id}/tags", middleware.Auth(server.RecordTagsGet))
	mux.HandleFunc("PUT /api/records/{record_id}/tags/{tag}", middleware.Auth(server.RecordTagsAttach))
	mux.HandleFunc("DELETE /api/records/{record_id}/tags/{tag}", middleware.Auth(server.RecordTagsDetach))

	httpServer := http.Server{
		Addr:    ":" + env.MustGet(envkeys.ServerPort),
		Handler: mux,
	}

	slog.Info("Server starting", "addr", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}
