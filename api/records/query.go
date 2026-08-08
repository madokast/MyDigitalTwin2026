package records

import (
	"context"
	"dt2026/httpx"
	"dt2026/lib"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type QueryCriteria struct {
	Page     *int64
	PageSize *int64
	Queries  []string
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
	criteria.Queries = q.GetOptionalStrings("q") // normalizeQueryCriteria 进一步检查

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
		total, err = queryToal(criteria, pool)
		if err != nil {
			return err
		}
	}

	totalPage := (total + *criteria.PageSize - 1) / *criteria.PageSize
	return &QueryRecordResponse{
		Ok:        true,
		Status:    http.StatusOK,
		Page:      *criteria.Page,
		PageSize:  *criteria.PageSize,
		Total:     total,
		TotalPage: totalPage,
		Records:   records,
	}

}

func queryToal(criteria *QueryCriteria, pool *pgxpool.Pool) (int64, *httpx.Error) {
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

	if enablePage {
		// ORDER BY id ASC
		sqlTail += " ORDER BY id ASC"

		// LIMIT
		args = append(args, *criteria.PageSize)
		sqlTail += fmt.Sprintf(" LIMIT $%d", len(args))

		// OFFSET
		args = append(args, (*criteria.Page-1)*(*criteria.PageSize))
		sqlTail += fmt.Sprintf(" OFFSET $%d", len(args))
	}

	return sqlTail, args
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

	// query 限制 10 个
	if len(criteria.Queries) > 10 {
		return httpx.NewBadRequestError(fmt.Sprintf(
			"too many queries, max 10, but got %d",
			len(criteria.Queries),
		))
	}

	// query 不能为空串
	for _, query := range criteria.Queries {
		if query == "" {
			return httpx.NewBadRequestError("query cannot be empty")
		}
	}

	return nil
}
