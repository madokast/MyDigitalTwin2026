package server

import (
	"dt2026/api/notify/qqbot"
	"dt2026/api/records"
	"dt2026/httpx"
	"net/http"
)

func (s *Server) RecordsAppend(w http.ResponseWriter, r *http.Request) {
	var record records.NewRecord
	if err := s.JsonUnmarshal(r, &record); err != nil {
		httpx.WriteJSON(w, err)
		return
	}
	pool, err := s.pool()
	if err != nil {
		httpx.WriteJSON(w, err)
		return
	}

	var bot *qqbot.Sender = nil
	if !s.testMode {
		bot, err = s.qqbot()
		if err != nil {
			httpx.WriteJSON(w, err)
			return
		}
	}

	httpx.WriteJSON(w, records.Append(&record, pool, bot))
}

func (s *Server) RecordsQuery(w http.ResponseWriter, r *http.Request) {
	q := httpx.Query(r.URL.Query())

	criteria, err := records.NewQueryCriteria(q)
	if err != nil {
		httpx.WriteJSON(w, err)
		return
	}

	pool, err := s.pool()
	if err != nil {
		httpx.WriteJSON(w, err)
		return
	}

	httpx.WriteJSON(w, records.Query(criteria, pool))
}
