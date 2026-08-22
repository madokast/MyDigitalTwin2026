package server

import (
	"dt2026/api/probe"
	"dt2026/httpx"
	"net/http"
)

func (s *Server) Time(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, probe.Time())
}
