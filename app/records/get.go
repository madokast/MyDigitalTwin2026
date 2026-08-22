package records

import (
	"context"
	"dt2026/httpx"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GetRecordsResponse struct {
	Ok     bool   `json:"ok"`
	Status int    `json:"status"`
	Record Record `json:"record"`
}

func (r GetRecordsResponse) GetStatus() int {
	return r.Status
}

var getByIDSQL = queryRecordSQL + " AND id = $1;"

func Get(recordID int64, pool *pgxpool.Pool) httpx.Response {
	var response = GetRecordsResponse{
		Ok:     true,
		Status: http.StatusOK,
	}

	var createAt time.Time
	var total int64
	err := pool.QueryRow(context.Background(), getByIDSQL, recordID).Scan(
		&response.Record.ID,
		&createAt,
		&response.Record.RawContent,
		&response.Record.ObjectiveContext,
		&response.Record.AIAnalysis,
		&response.Record.Tags,
		&total,
	)
	_ = total // unused

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return httpx.NewNotFoundError(fmt.Sprintf(
				"record %d not found", recordID,
			))
		}

		return httpx.NewInternalServerError(fmt.Sprintf(
			"failed to query %s with %v: %s",
			getByIDSQL, recordID, err.Error(),
		))
	}

	response.Record.CreatedAt = JSONTime(createAt)

	return response
}
