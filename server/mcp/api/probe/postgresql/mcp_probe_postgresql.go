package postgresql

import (
	"context"
	"dt2026/api/probe"
	"dt2026/httpx"
	"dt2026/server/http"
)

const Name = "probe_postgresql"

const Description = "Probe whether PostgreSQL is reachable. Opens a fresh connection, runs SELECT NOW()::text, and returns connection/query latency in milliseconds plus the database clock. Do not use this to query records or write data."

type Input struct {
}

type Output = probe.ProbePostgreSQLResponse

func ProbePostgreSQL(_ *http.Server, _ Input) (Output, *httpx.Error) {
	res := probe.ProbePostgreSQL(context.Background())
	if httpErr, ok := res.(*httpx.Error); ok {
		var zero Output
		return zero, httpErr
	}
	return res.(probe.ProbePostgreSQLResponse), nil
}
