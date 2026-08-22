package server

import (
	"dt2026/app/probe"
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
	q := httpx.QueryParams(r.URL.Query())

	message, err := q.GetOptionalSingleString("message")
	if err != nil {
		httpx.WriteJSON(w, err)
		return
	}

	bot, err := s.QQBot()
	if err != nil {
		httpx.WriteJSON(w, err)
		return
	}
	httpx.WriteJSON(w, probe.ProbeQQBot(message, bot))
}
