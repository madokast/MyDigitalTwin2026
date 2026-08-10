package server

import (
	"dt2026/api/notify/qqbot"
	"dt2026/api/records"
	"dt2026/httpx"
	"net/http"
)

func (s *Server) RecordsPost(w http.ResponseWriter, r *http.Request) {
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

	httpx.WriteJSON(w, records.Post(&record, pool, bot))
}

func (s *Server) RecordsQuery(w http.ResponseWriter, r *http.Request) {
	q := httpx.QueryParams(r.URL.Query())

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

func (s *Server) RecordGet(w http.ResponseWriter, r *http.Request) {
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

	httpx.WriteJSON(w, records.Get(recordID, pool))
}

func (s *Server) RecordsExport(w http.ResponseWriter, r *http.Request) {
	q := httpx.QueryParams(r.URL.Query())

	request, err := records.NewExportRequest(q)
	if err != nil {
		httpx.WriteJSON(w, err)
		return
	}

	pool, err := s.pool()
	if err != nil {
		httpx.WriteJSON(w, err)
		return
	}

	err = records.Export(request, w, pool)
	if err != nil {
		httpx.WriteJSON(w, err)
		return
	}

	// 无 err，则导出数据已流式输出
}
