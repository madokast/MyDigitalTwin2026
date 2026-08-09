package tags

import (
	"context"
	"dt2026/httpx"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GetTagResponse struct {
	Ok     bool     `json:"ok"`
	Status int      `json:"status"`
	Tags   []string `json:"tags"`
}

func (r GetTagResponse) GetStatus() int {
	return r.Status
}

func Get(recordID int64, pool *pgxpool.Pool) httpx.Response {

	var response = GetTagResponse{
		Ok:     true,
		Status: http.StatusOK,
	}

	err := pool.QueryRow(
		context.Background(),
		getRecordTagsSQL,
		recordID,
	).Scan(&response.Tags)

	if err != nil {
		if err == pgx.ErrNoRows {
			return httpx.NewNotFoundError(fmt.Sprintf(
				"record %d not found", recordID,
			))
		}

		return httpx.NewInternalServerError(fmt.Sprintf(
			"failed to query %s with %v: %s",
			getRecordTagsSQL, recordID, err.Error(),
		))
	}

	return response
}
