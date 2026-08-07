package main

import (
	"dt2026/api/envkeys"
	"dt2026/api/middleware"
	"dt2026/api/server"
	"dt2026/env"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	slog.Info("Starting")

	server := server.NewServer()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/probe/health", middleware.Auth(server.ProbeHealth))
	mux.HandleFunc("GET /api/probe/bad-json", middleware.Auth(server.ProbeBadJson))
	mux.HandleFunc("GET /api/probe/postgresql", middleware.Auth(server.ProbePostgreSQL))
	mux.HandleFunc("GET /api/probe/qqbot", middleware.Auth(server.ProbeQQBot))
	mux.HandleFunc("POST /api/records", middleware.Auth(server.RecordsAppend))

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
