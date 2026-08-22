package records

import (
	"context"
	"dt2026/httpx"
	"dt2026/lib/optional"
	"encoding/json/v2"
	"fmt"
	"log/slog"
	"net/http"
	"time"

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
	n, err := w.w.Write(p)
	if err != nil {
		slog.Error("failed to write response", "err", err)
	}
	return n, err
}

func Export(r *ExportRequest, w http.ResponseWriter,
	pool *pgxpool.Pool, maxExportSize optional.Optional[int64]) *httpx.Error {

	httpErr := validateExportSize(r, pool, maxExportSize)
	if httpErr != nil {
		return httpErr
	}
	// validate 和后续导出不开启事务，允许中途存在新增，可能导致数目稍微超过限制

	sqlTail, args := makeQuerySQLTail(&QueryCriteria{
		From: r.From,
		To:   r.To,
	}, DisablePage)
	sql := fmt.Sprintf(exportRecordSQL, sqlTail)

	rows, err := pool.Query(context.Background(), sql, args...)
	if err != nil {
		return httpx.NewInternalServerError(fmt.Sprintf(
			"failed to query %s with %v: %s",
			sql, args, err.Error(),
		))
	}
	defer rows.Close()

	cw := checkableWriter{w: w}
	for rows.Next() {
		var record Record
		var createdAt time.Time

		err := rows.Scan(
			&record.ID,
			&createdAt,
			&record.RawContent,
			&record.ObjectiveContext,
			&record.AIAnalysis,
			&record.Tags,
		)

		if err != nil {
			msg := fmt.Sprintf(
				"failed to scan row in querying %s with %v: %s",
				sql, args, err.Error(),
			)
			if cw.started {
				// 传输到一半，无法改写响应头
				_, _ = cw.Write([]byte("export failed!\n"))
				_, _ = cw.Write([]byte(msg))
				return nil
			} else {
				return httpx.NewInternalServerError(msg)
			}
		}

		record.CreatedAt = JSONTime(createdAt)

		data, err := json.Marshal(record)
		if err != nil {
			msg := fmt.Sprintf(
				"failed to json marshal record %v: %s",
				record, err.Error(),
			)
			if cw.started {
				// 传输到一半，无法改写响应头
				_, _ = cw.Write([]byte("export failed!\n"))
				_, _ = cw.Write([]byte(msg))
				return nil
			} else {
				return httpx.NewInternalServerError(msg)
			}
		}
		_, _ = cw.Write(data)
		_, _ = cw.Write([]byte{'\n'})
	}

	if err := rows.Err(); err != nil {
		msg := fmt.Sprintf(
			"failed to iterate rows in querying %s with %v: %s",
			sql, args, err.Error(),
		)
		if cw.started {
			// 传输到一半，无法改写响应头
			_, _ = cw.Write([]byte("export failed!\n"))
			_, _ = cw.Write([]byte(msg))
			return nil
		} else {
			return httpx.NewInternalServerError(msg)
		}
	}

	return nil
}

func validateExportSize(r *ExportRequest, pool *pgxpool.Pool,
	maxExportSize optional.Optional[int64]) *httpx.Error {

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
			message += ". Provide a specific [from, to) time range."
		} else {
			message += ". Try a shorter [from, to) time range."
		}

		return httpx.NewBadRequestError(message)
	}

	return nil
}
