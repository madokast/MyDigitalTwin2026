package tags

import (
	"context"
	"dt2026/httpx"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TagCount struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

type GetAllTagsResponse struct {
	Ok        httpx.TrueType `json:"ok"`
	Status    int        `json:"status"`
	TagCounts []TagCount `json:"tag_counts"`
}

func (r GetAllTagsResponse) GetStatus() int {
	return r.Status
}

func GetAll(pool *pgxpool.Pool) httpx.Response {

	var response = GetAllTagsResponse{
		Ok:     httpx.True,
		Status: http.StatusOK,
	}

	const sql = getAllTagCountSQL
	rows, err := pool.Query(
		context.Background(),
		sql,
	)
	if err != nil {
		return httpx.NewInternalServerError(fmt.Sprintf(
			"failed to query %s: %s",
			sql, err.Error(),
		))
	}
	defer rows.Close()

	for rows.Next() {
		var tc TagCount

		err := rows.Scan(
			&tc.Tag,
			&tc.Count,
		)

		if err != nil {
			return httpx.NewInternalServerError(fmt.Sprintf(
				"failed to scan row in querying %s: %s",
				sql, err.Error(),
			))
		}

		response.TagCounts = append(response.TagCounts, tc)
	}

	if err := rows.Err(); err != nil {
		return httpx.NewInternalServerError(fmt.Sprintf(
			"failed to iterate rows in querying %s: %s",
			sql, err.Error(),
		))
	}

	return response
}
