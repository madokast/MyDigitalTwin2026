package tags

import (
	"context"
	"dt2026/httpx"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DetachTagResponse struct {
	Ok       httpx.TrueType `json:"ok" jsonschema:"Whether the request succeeded"`
	Status   int            `json:"status" jsonschema:"HTTP status code"`
	Detached bool           `json:"detached" jsonschema:"Whether this call removed the tag; false if it was already absent"`
	Changed  bool           `json:"changed" jsonschema:"Whether this call removed the tag; false if it was already absent"`
	Tags     []string       `json:"tags" jsonschema:"Tags on this record after the detach"`
}

func (r DetachTagResponse) GetStatus() int {
	return r.Status
}

func Detach(recordID int64, tag string, pool *pgxpool.Pool) httpx.Response {
	var httpErr *httpx.Error
	tag, httpErr = NormalizeTag(tag)
	if httpErr != nil {
		return httpErr
	}

	var response = DetachTagResponse{
		Ok:       httpx.True,
		Status:   http.StatusOK,
		Detached: true,
		Tags:     []string{},
	}

	var detachResult int
	err := pool.QueryRow(
		context.Background(),
		detachTagSQL,
		recordID,
		tag,
	).Scan(&detachResult, &response.Tags)

	if err != nil {
		return httpx.NewInternalServerError(fmt.Sprintf(
			"failed to execute %s with %v: %s",
			detachTagSQL, []any{recordID, tag}, err.Error(),
		))
	}

	switch detachResult {
	case detachResultTagDetached:
		response.Changed = true
	case detachResultRecordNotFound:
		return httpx.NewNotFoundError(fmt.Sprintf(
			"record %d not found", recordID,
		))
	case detachResultTagNotExists:
		response.Detached = false
	default:
		return httpx.NewInternalServerError(fmt.Sprintf(
			"unknown detach result %d while executing %s with %v",
			detachResult, detachTagSQL, []any{recordID, tag},
		))
	}

	if response.Tags == nil {
		response.Tags = []string{}
	}

	return response
}
