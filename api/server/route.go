package server

import (
	"dt2026/api/probe"
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
