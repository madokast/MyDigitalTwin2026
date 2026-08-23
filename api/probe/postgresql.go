package probe

import (
	"context"
	"dt2026/api/envkeys"
	"dt2026/env"
	"dt2026/httpx"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

type ProbePostgreSQLResponse struct {
	Ok                 httpx.TrueType `json:"ok" jsonschema:"Whether the request succeeded"`
	Status             int    `json:"status" jsonschema:"HTTP status code"`
	ConnectionTimeMs   int64  `json:"connection_time_ms" jsonschema:"Time to open a new database connection, in milliseconds"`
	QueryTimeMs        int64  `json:"query_time_ms" jsonschema:"Time to run SELECT NOW()::text, in milliseconds"`
	QueryNowTextResult string `json:"query_now_text_result" jsonschema:"Database clock as returned by SELECT NOW()::text"`
}

func (r ProbePostgreSQLResponse) GetStatus() int {
	return r.Status
}

func ProbePostgreSQL(ctx context.Context) httpx.Response {
	var response = ProbePostgreSQLResponse{
		Ok:     httpx.True,
		Status: http.StatusOK,
	}

	start := time.Now()

	url, ok := env.Get(envkeys.DatabaseUrl)
	if !ok {
		return httpx.NewInternalServerError("database url not set")
	}

	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return httpx.NewInternalServerError("failed to connect to database: " + err.Error())
	}
	response.ConnectionTimeMs = time.Since(start).Milliseconds()
	start = time.Now()

	defer func() {
		_ = conn.Close(ctx)
	}()

	err = conn.QueryRow(ctx, "SELECT NOW()::text").Scan(&response.QueryNowTextResult)
	if err != nil {
		return httpx.NewInternalServerError("failed to query database: " + err.Error())
	}
	response.QueryTimeMs = time.Since(start).Milliseconds()
	return response
}
