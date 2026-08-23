package tags

import (
	"context"
	"dt2026/httpx"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AttachTagResponse struct {
	Ok       httpx.TrueType `json:"ok" jsonschema:"Whether the request succeeded"`
	Status   int            `json:"status" jsonschema:"HTTP status code"`
	Attached bool           `json:"attached" jsonschema:"Whether the tag is present on the record after this call"`
	Changed  bool           `json:"changed" jsonschema:"Whether this call added the tag; false if it was already present"`
	Tags     []string       `json:"tags" jsonschema:"Tags on this record after the attach"`
}

func (r AttachTagResponse) GetStatus() int {
	return r.Status
}

func Attach(recordID int64, tag string, pool *pgxpool.Pool) httpx.Response {
	var httpErr *httpx.Error
	tag, httpErr = NormalizeTag(tag)
	if httpErr != nil {
		return httpErr
	}

	var response = AttachTagResponse{
		Ok:       httpx.True,
		Attached: true,
		Tags:     []string{},
	}

	var attachResult int
	err := pool.QueryRow(
		context.Background(),
		attachTagSQL,
		recordID,
		tag,
	).Scan(&attachResult, &response.Tags)

	if err != nil {
		return httpx.NewInternalServerError(fmt.Sprintf(
			"failed to execute %s with %v: %s",
			attachTagSQL, []any{recordID, tag}, err.Error(),
		))
	}

	switch attachResult {
	case attachResultTagAttached:
		response.Status = http.StatusCreated
		response.Changed = true
	case attachResultRecordNotFound:
		return httpx.NewNotFoundError(fmt.Sprintf(
			"record %d not found", recordID,
		))
	case attachResultTagAlreadyExists:
		response.Status = http.StatusOK
		response.Changed = false
	default:
		return httpx.NewInternalServerError(fmt.Sprintf(
			"unknown attach result %d while executing %s with %v",
			attachResult, attachTagSQL, []any{recordID, tag},
		))
	}

	if response.Tags == nil {
		response.Tags = []string{}
	}

	return response
}
