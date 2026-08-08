package records

import (
	"context"
	"dt2026/httpx"
	"dt2026/lib"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type QueryCriteria struct {
	Page     *int64
	PageSize *int64
}

func NewQueryCriteria(q httpx.Query) (*QueryCriteria, *httpx.Error) {
	var criteria QueryCriteria
	var err *httpx.Error

	criteria.Page, err = q.GetOptionalSingleInt64("page")
	if err != nil {
		return nil, err
	}
	criteria.PageSize, err = q.GetOptionalSingleInt64("page_size")
	if err != nil {
		return nil, err
	}

	return &criteria, nil
}

type QueryRecordResponse struct {
	Ok       bool     `json:"ok"`
	Status   int      `json:"status"`
	Page     int64    `json:"page"`
	PageSize int64    `json:"page_size"`
	Records  []Record `json:"records"`
}

func (r QueryRecordResponse) GetStatus() int {
	return r.Status
}

func Query(criteria *QueryCriteria, pool *pgxpool.Pool) httpx.Response {
	if err := normalizeQueryCriteria(criteria); err != nil {
		return err
	}

	sql, args := makeQuerySQL(criteria)

	rows, err := pool.Query(context.Background(), sql, args...)
	if err != nil {
		return httpx.NewError(http.StatusInternalServerError, fmt.Sprintf(
			"failed to query %s with %v: %s",
			sql, args, err.Error(),
		))
	}
	defer rows.Close()

	var records []Record

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
			return httpx.NewError(http.StatusInternalServerError, fmt.Sprintf(
				"failed to scan row in querying %s with %v: %s",
				sql, args, err.Error(),
			))
		}

		record.CreatedAt = JSONTime(createdAt)
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return httpx.NewError(http.StatusInternalServerError, fmt.Sprintf(
			"failed to iterate rows in querying %s with %v: %s",
			sql, args, err.Error(),
		))
	}

	return &QueryRecordResponse{
		Ok:       true,
		Status:   http.StatusOK,
		Page:     *criteria.Page,
		PageSize: *criteria.PageSize,
		Records:  records,
	}

}

func makeQuerySQL(criteria *QueryCriteria) (sql string, args []any) {

	sql = queryRecordSQL

	// ORDER BY id ASC
	sql += " ORDER BY id ASC"

	// LIMIT
	args = append(args, *criteria.PageSize)
	sql += fmt.Sprintf(" LIMIT $%d", len(args))

	// OFFSET
	args = append(args, (*criteria.Page-1)*(*criteria.PageSize))
	sql += fmt.Sprintf(" OFFSET $%d", len(args))

	return sql, args
}

func normalizeQueryCriteria(criteria *QueryCriteria) *httpx.Error {
	if criteria.Page == nil {
		criteria.Page = new(int64)
		*criteria.Page = 1
	}
	if *criteria.Page <= 0 {
		return httpx.NewBadRequestError(fmt.Sprintf(
			"page must be greater than 0, but got %d",
			*criteria.Page,
		))
	}

	if criteria.PageSize == nil {
		criteria.PageSize = new(int64)
		*criteria.PageSize = 100
	}
	if *criteria.PageSize <= 0 || *criteria.PageSize > 1000 {
		return httpx.NewBadRequestError(fmt.Sprintf(
			"page_size must be between 1 and 1000, but got %d",
			*criteria.PageSize,
		))
	}

	// OFFSET 不能溢出
	if lib.MulOverflow(*criteria.Page-1, *criteria.PageSize) {
		return httpx.NewBadRequestError(fmt.Sprintf(
			"page and page_size overflow, page=%d, page_size=%d",
			*criteria.Page, *criteria.PageSize,
		))
	}

	return nil
}
