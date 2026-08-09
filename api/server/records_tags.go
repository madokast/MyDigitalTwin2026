package server

import (
	"dt2026/api/records/tags"
	"dt2026/httpx"
	"net/http"
)

func (s *Server) RecordTagsGet(w http.ResponseWriter, r *http.Request) {
	pv := httpx.PathParams{Request: r}

	recordID, err := pv.Int64("record_id")
	if err != nil {
		httpx.WriteJSON(w, err)
		return
	}

	pool, err := s.pool()
	if err != nil {
		httpx.WriteJSON(w, err)
		return
	}

	httpx.WriteJSON(w, tags.Get(recordID, pool))
}

func (s *Server) RecordTagsAttach(w http.ResponseWriter, r *http.Request) {
	pv := httpx.PathParams{Request: r}

	recordID, err := pv.Int64("record_id")
	if err != nil {
		httpx.WriteJSON(w, err)
		return
	}

	tag, err := pv.String("tag")
	if err != nil {
		httpx.WriteJSON(w, err)
		return
	}

	pool, err := s.pool()
	if err != nil {
		httpx.WriteJSON(w, err)
		return
	}

	httpx.WriteJSON(w, tags.Attach(recordID, tag, pool))
}
