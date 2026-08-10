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

const defaultMaxExportSize = 10000

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
	pool *pgxpool.Pool, maxExportSize optional.Optional[int64]) *httpx.Error {

	total, httpErr := queryTotal(&QueryCriteria{
		From: r.From,
		To:   r.To,
	}, pool)
	if httpErr != nil {
		return httpErr
	}

	maxSize := maxExportSize.Or(defaultMaxExportSize)

	if total > maxSize {
		message := fmt.Sprintf(
			"too many records to export: total=%d, limit=%d",
			total, maxSize,
		)
		if r.From.Absent() || r.To.Absent() {
			message += ". Provide a specific from/to time range."
		} else {
			message += ". Try a shorter from/to time range."
		}

		return httpx.NewBadRequestError(message)
	} else if total == 0 {
		// OK. 保存进一个空文件
		w.Header().Set("Content-Type", "application/x-ndjson")
		w.WriteHeader(http.StatusOK)
		w.Write(nil)
		return nil
	}

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

	sql := fmt.Sprintf(exportRecordSQL, predicates, maxSize)

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
