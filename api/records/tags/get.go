package tags

import (
	"context"
	"dt2026/httpx"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GetTagResponse struct {
	Ok     httpx.TrueType `json:"ok" jsonschema:"Whether the request succeeded"`
	Status int            `json:"status" jsonschema:"HTTP status code"`
	Tags   []string       `json:"tags" jsonschema:"Tags on this record as a string array; empty if none"`
}

func (r GetTagResponse) GetStatus() int {
	return r.Status
}

func Get(recordID int64, pool *pgxpool.Pool) httpx.Response {

	var response = GetTagResponse{
		Ok:     httpx.True,
		Status: http.StatusOK,
		Tags:   []string{},
	}

	err := pool.QueryRow(
		context.Background(),
		getRecordTagsSQL,
		recordID,
	).Scan(&response.Tags)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return httpx.NewNotFoundError(fmt.Sprintf(
				"record %d not found", recordID,
			))
		}

		return httpx.NewInternalServerError(fmt.Sprintf(
			"failed to query %s with %v: %s",
			getRecordTagsSQL, recordID, err.Error(),
		))
	}

	if response.Tags == nil {
		response.Tags = []string{}
	}

	return response
}
