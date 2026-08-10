package records

import (
	"context"
	"dt2026/httpx"
	"dt2026/lib/optional"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ExportRequest struct {
	From optional.Optional[JSONTime]
	To   optional.Optional[JSONTime]
}

func NewExportRequest(q httpx.QueryParams) (*ExportRequest, *httpx.Error) {
	var request ExportRequest

	from, httpErr := q.GetOptionalSingleString("from")
	if httpErr != nil {
		return nil, httpErr
	}
	request.From, httpErr = ParseOptionalJSONTime(from)
	if httpErr != nil {
		return nil, httpErr
	}

	to, httpErr := q.GetOptionalSingleString("to")
	if httpErr != nil {
		return nil, httpErr
	}
	request.To, httpErr = ParseOptionalJSONTime(to)
	if httpErr != nil {
		return nil, httpErr
	}

	return &request, nil
}

type checkableWriter struct {
	w       http.ResponseWriter
	started bool
}

func (w *checkableWriter) Write(p []byte) (int, error) {
	if !w.started {
		w.started = true
		w.w.Header().Set("Content-Type", "application/x-ndjson")
		w.w.WriteHeader(http.StatusOK)
	}
	return w.w.Write(p)
}

func Export(r *ExportRequest, w http.ResponseWriter,
	pool *pgxpool.Pool) *httpx.Error {

	var predicates string

	// PostgreSQL timestamptz 精度到微秒
	const RFC3339Micro = "2006-01-02T15:04:05.999999Z07:00"
	if from, ok := r.From.Get(); ok {
		predicates += fmt.Sprintf(
			" AND created_at >= '%s'::timestamptz",
			from.GoTime().UTC().Format(RFC3339Micro),
		)
	}
	if to, ok := r.To.Get(); ok {
		predicates += fmt.Sprintf(
			" AND created_at < '%s'::timestamptz",
			to.GoTime().UTC().Format(RFC3339Micro),
		)
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return httpx.NewInternalServerError(fmt.Sprintf(
			"failed to acquire pgx connection from pool: %s",
			err.Error(),
		))
	}
	defer conn.Release()

	sql := fmt.Sprintf(exportRecordSQL, predicates)

	cw := checkableWriter{w: w}
	_, err = conn.Conn().PgConn().CopyTo(
		context.Background(),
		&cw,
		sql,
	)

	if err != nil {
		slog.Error(
			"failed to export records while executing copy-to",
			"sql", sql, "err", err,
		)
		if !cw.started {
			return httpx.NewInternalServerError(fmt.Sprintf(
				"failed to export records while executing %s: %s",
				sql, err.Error(),
			))
		} else {
			// 输出到到一半时 PostgreSQL 出错，无法改写响应头
			// 符合 HTTP streaming 的正常语义
		}
	}

	return nil
}
