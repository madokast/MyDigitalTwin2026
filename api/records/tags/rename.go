package tags

import (
	"context"
	"dt2026/httpx"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RenameTagRequest struct {
	NewTag string `json:"new_tag" jsonschema:"New name for the tag; required, trimmed, must not be blank"`
}

type RenameTagResponse struct {
	Ok      httpx.TrueType `json:"ok" jsonschema:"Whether the request succeeded"`
	Status  int            `json:"status" jsonschema:"HTTP status code"`
	Renamed bool           `json:"renamed" jsonschema:"Whether the tag bears the new name on the record after this call"`
	Changed bool           `json:"changed" jsonschema:"Whether this call renamed the tag; false if it already bore the new name or was absent"`
	Tags    []string       `json:"tags" jsonschema:"Tags on this record after the rename"`
}

func (r RenameTagResponse) GetStatus() int {
	return r.Status
}

func Rename(recordID int64, tag string, newTag string, pool *pgxpool.Pool) httpx.Response {
	var httpErr *httpx.Error
	tag, httpErr = NormalizeTag(tag)
	if httpErr != nil {
		return httpErr
	}

	newTag, httpErr = NormalizeTag(newTag)
	if httpErr != nil {
		return httpErr
	}

	var response = RenameTagResponse{
		Ok:     httpx.True,
		Status: http.StatusOK,
		Tags:   []string{},
	}

	var renameResult int
	err := pool.QueryRow(
		context.Background(),
		renameTagSQL,
		recordID,
		tag,
		newTag,
	).Scan(&renameResult, &response.Tags)

	if err != nil {
		return httpx.NewInternalServerError(fmt.Sprintf(
			"failed to execute %s with %v: %s",
			renameTagSQL, []any{recordID, tag, newTag}, err.Error(),
		))
	}

	switch renameResult {
	case renameResultTagRenamed:
		response.Renamed = true
		response.Changed = true
	case renameResultRecordNotFound:
		return httpx.NewNotFoundError(fmt.Sprintf(
			"record %d not found", recordID,
		))
	case renameResultTagNotExists:
		response.Renamed = false
	case renameResultNewTagSameAsOld:
		response.Renamed = true
	case renameResultNewTagExists:
		return httpx.NewError(http.StatusConflict, fmt.Sprintf(
			"tag '%s' already exists on record %d", newTag, recordID,
		))
	default:
		return httpx.NewInternalServerError(fmt.Sprintf(
			"unknown rename result %d while executing %s with %v",
			renameResult, renameTagSQL, []any{recordID, tag, newTag},
		))
	}

	if response.Tags == nil {
		response.Tags = []string{}
	}

	return response
}
