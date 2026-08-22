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
		w.writeHeader()
	}
	n, err := w.w.Write(p)
	if err != nil {
		slog.Error("failed to write response", "err", err)
	}
	return n, err
}

func (w *checkableWriter) writeHeader() {
	w.w.Header().Set("Content-Type", "application/x-ndjson")
	w.w.WriteHeader(http.StatusOK)
}

func (w *checkableWriter) sendError(err *httpx.Error) *httpx.Error {
	if w.started {
		// 传输到一半，无法改写响应头，错误输出到 data 中
		httpx.WriteJSON(w.w, err)
		return nil
	}
	// 传输未开始，外部负责写错误
	return err
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
	exportCount := 0
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
			return cw.sendError(httpx.NewInternalServerError(fmt.Sprintf(
				"failed to export records while scanning row in querying %s with %v: %s",
				sql, args, err.Error(),
			)))
		}

		record.CreatedAt = JSONTime(createdAt)

		data, err := json.Marshal(record)
		if err != nil {
			return cw.sendError(httpx.NewInternalServerError(fmt.Sprintf(
				"failed to export records while json marshaling record %v: %s",
				record, err.Error(),
			)))
		}
		_, _ = cw.Write(data)
		_, _ = cw.Write([]byte{'\n'})
		exportCount++
	}

	if err := rows.Err(); err != nil {
		return cw.sendError(httpx.NewInternalServerError(fmt.Sprintf(
			"failed to export records while iterating rows in querying %s with %v: %s",
			sql, args, err.Error(),
		)))
	}

	if exportCount == 0 {
		cw.writeHeader()
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
