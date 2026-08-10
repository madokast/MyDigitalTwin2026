package records

import (
	"context"
	"dt2026/api/records/tags"
	"dt2026/httpx"
	"dt2026/lib"
	"fmt"
	"net/http"
	"strings"
	"time"

	"dt2026/lib/optional"

	"github.com/jackc/pgx/v5/pgxpool"
)

type QueryCriteria struct {
	Page     optional.Optional[int64]
	PageSize optional.Optional[int64]
	Queries  []string
	Tags     []string
}

const (
	defaultPage     int64 = 1
	defaultPageSize int64 = 100
)

func NewQueryCriteria(q httpx.QueryParams) (*QueryCriteria, *httpx.Error) {
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

	criteria.Queries = q.GetOptionalStrings("q") // normalizeQueryCriteria 进一步检查

	criteria.Tags = q.GetOptionalStrings("tag") // normalizeQueryCriteria 进一步检查

	return &criteria, nil
}

type QueryRecordResponse struct {
	Ok        bool     `json:"ok"`
	Status    int      `json:"status"`
	Page      int64    `json:"page"`
	PageSize  int64    `json:"page_size"`
	Total     int64    `json:"total"`
	TotalPage int64    `json:"total_page"`
	Records   []Record `json:"records"`
}

func (r QueryRecordResponse) GetStatus() int {
	return r.Status
}

func Query(criteria *QueryCriteria, pool *pgxpool.Pool) httpx.Response {
	if err := normalizeQueryCriteria(criteria); err != nil {
		return err
	}

	sqlTail, args := makeQuerySQLTail(criteria, true)
	sql := queryRecordSQL + sqlTail

	rows, err := pool.Query(context.Background(), sql, args...)
	if err != nil {
		return httpx.NewError(http.StatusInternalServerError, fmt.Sprintf(
			"failed to query %s with %v: %s",
			sql, args, err.Error(),
		))
	}
	defer rows.Close()

	var records []Record
	var total int64

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
			&total,
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

	if len(records) == 0 { // 没有扫描到行，则 total 拿不到，单独用 COUNT 查询
		var err *httpx.Error
		total, err = queryTotal(criteria, pool)
		if err != nil {
			return err
		}
	}

	page := criteria.Page.Or(defaultPage)
	pageSize := criteria.PageSize.Or(defaultPageSize)

	totalPage := (total + pageSize - 1) / pageSize
	return QueryRecordResponse{
		Ok:        true,
		Status:    http.StatusOK,
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: totalPage,
		Records:   records,
	}

}

func queryTotal(criteria *QueryCriteria, pool *pgxpool.Pool) (int64, *httpx.Error) {
	sqlTail, args := makeQuerySQLTail(criteria, false)
	sql := countRecordSQL + sqlTail

	var total int64
	err := pool.QueryRow(
		context.Background(),
		sql,
		args...,
	).Scan(&total)

	if err != nil {
		return 0, httpx.NewError(http.StatusInternalServerError, fmt.Sprintf(
			"failed to query %s with %v: %s",
			sql, args, err.Error(),
		))
	}

	return total, nil
}

// makeQuerySQLTail 构造 SQL 语句
// enablePage 带上 ORDER、LIMIT、OFFSET 信息
func makeQuerySQLTail(criteria *QueryCriteria, enablePage bool) (sqlTail string, args []any) {
	// querys
	for _, query := range criteria.Queries {
		var predicates []string
		queryWord := "%" + lib.EscapeSQLLike(query) + "%"
		args = append(args, queryWord)
		for _, column := range textColumns { // TEXT 字段模糊搜索
			predicates = append(predicates,
				fmt.Sprintf("%s ILIKE $%d", column, len(args)))
		}
		// tags 字段模糊搜索
		predicates = append(predicates,
			fmt.Sprintf("EXISTS (SELECT 1 FROM unnest(tags) AS tag WHERE tag ILIKE $%d)", len(args)))

		// 组装 SQL
		fragment := strings.Join(predicates, " OR ")
		sqlTail += fmt.Sprintf(" AND (%s)", fragment)
	}

	// tags
	var tagPlaceholders []string
	for _, tag := range criteria.Tags {
		args = append(args, tag)
		tagPlaceholders = append(tagPlaceholders, fmt.Sprintf("$%d", len(args)))
	}
	if len(tagPlaceholders) > 0 {
		sqlTail += fmt.Sprintf(
			" AND tags @> ARRAY[%s]::TEXT[]",
			strings.Join(tagPlaceholders, ", "),
		)
	}

	if enablePage {
		// ORDER BY id ASC
		sqlTail += " ORDER BY id ASC"

		page := criteria.Page.Or(defaultPage)
		pageSize := criteria.PageSize.Or(defaultPageSize)

		// LIMIT
		args = append(args, pageSize)
		sqlTail += fmt.Sprintf(" LIMIT $%d", len(args))

		// OFFSET
		args = append(args, (page-1)*pageSize)
		sqlTail += fmt.Sprintf(" OFFSET $%d", len(args))
	}

	return sqlTail, args
}

func normalizeQueryCriteria(criteria *QueryCriteria) *httpx.Error {
	page := criteria.Page.Or(defaultPage)
	pageSize := criteria.PageSize.Or(defaultPageSize)

	if page <= 0 {
		return httpx.NewBadRequestError(fmt.Sprintf(
			"page must be greater than 0, but got: '%d'",
			page,
		))
	}

	if pageSize <= 0 || pageSize > 1000 {
		return httpx.NewBadRequestError(fmt.Sprintf(
			"page_size must be between 1 and 1000, but got: '%d'",
			pageSize,
		))
	}

	// OFFSET 不能溢出
	if lib.MulOverflow(page-1, pageSize) {
		return httpx.NewBadRequestError(fmt.Sprintf(
			"page and page_size overflow, page=%d, page_size=%d",
			page, pageSize,
		))
	}

	// queries 限制 10 个
	if len(criteria.Queries) > 10 {
		return httpx.NewBadRequestError(fmt.Sprintf(
			"too many queries, max 10, but got %d queries: (%s)",
			len(criteria.Queries), lib.SliceToString(criteria.Queries),
		))
	}

	// query 不能为空串
	for _, query := range criteria.Queries {
		if query == "" {
			return httpx.NewBadRequestError("query cannot be empty")
		}
	}

	// tags 限制 10 个
	if len(criteria.Tags) > 10 {
		return httpx.NewBadRequestError(fmt.Sprintf(
			"too many tags, max 10, but got %d tags: (%s)",
			len(criteria.Tags), lib.SliceToString(criteria.Tags),
		))
	}

	// tag normalize
	var err *httpx.Error
	criteria.Tags, err = tags.NormalizeTags(criteria.Tags)
	if err != nil {
		return err
	}

	return nil
}
