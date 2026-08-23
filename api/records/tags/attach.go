package tags

import (
	"context"
	"dt2026/httpx"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AttachTagResponse struct {
	Ok       httpx.TrueType `json:"ok"`
	Status   int      `json:"status"`
	Attached bool     `json:"attached"`
	Changed  bool     `json:"changed"`
	Tags     []string `json:"tags"`
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

	return response
}
