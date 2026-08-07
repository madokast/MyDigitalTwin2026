package server

import (
	"dt2026/api/probe"
	"dt2026/api/records"
	"dt2026/httpx"
	"net/http"
)

func (s *Server) ProbeHealth(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, probe.Health())
}

func (s *Server) ProbeBadJson(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, probe.BadJSON())
}

func (s *Server) ProbePostgreSQL(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, probe.ProbePostgreSQL(r.Context()))
}

func (s *Server) ProbeQQBot(w http.ResponseWriter, r *http.Request) {
	bot, err := s.qqbot()
	if err != nil {
		httpx.WriteJSON(w, err)
		return
	}
	httpx.WriteJSON(w, probe.ProbeQQBot(bot))
}

func (s *Server) RecordsAppend(w http.ResponseWriter, r *http.Request) {
	var record records.NewRecord
	if err := s.JsonUnmarshal(r, &record); err != nil {
		httpx.WriteJSON(w, err)
		return
	}
	pool, err := s.db()
	if err != nil {
		httpx.WriteJSON(w, err)
		return
	}

	httpx.WriteJSON(w, records.Append(&record, pool))
}
