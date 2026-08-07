package probe

import (
	"dt2026/api/envkeys"
	"dt2026/api/httpx"
	"dt2026/env"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

type ProbePostgreSQLResponse struct {
	Ok                 bool   `json:"ok"`
	Status             int    `json:"status"`
	ConnectionTimeMs   int64  `json:"connection_time_ms"`
	QueryTimeMs        int64  `json:"query_time_ms"`
	QueryNowTextResult string `json:"query_now_text_result"`
}

func (r ProbePostgreSQLResponse) IsOk() bool {
	return r.Ok
}

func (r ProbePostgreSQLResponse) GetStatus() int {
	return r.Status
}

func ProbePostgreSQL(w http.ResponseWriter, r *http.Request) {
	var response = ProbePostgreSQLResponse{
		Ok:     true,
		Status: http.StatusOK,
	}

	start := time.Now()

	url, ok := env.Get(envkeys.DatabaseUrl)
	if !ok {
		httpx.WriteJSON(w, httpx.NewInternalServerError("database url not set"))
		return
	}

	conn, err := pgx.Connect(r.Context(), url)
	if err != nil {
		httpx.WriteJSON(w, httpx.NewInternalServerError("failed to connect to database: "+err.Error()))
		return
	}
	response.ConnectionTimeMs = time.Since(start).Milliseconds()
	start = time.Now()

	defer func() {
		_ = conn.Close(r.Context())
	}()

	err = conn.QueryRow(r.Context(), "SELECT NOW()::text").Scan(&response.QueryNowTextResult)
	if err != nil {
		httpx.WriteJSON(w, httpx.NewInternalServerError("failed to query database: "+err.Error()))
		return
	}
	response.QueryTimeMs = time.Since(start).Milliseconds()
	httpx.WriteJSON(w, response)
}
