package query

import (
	"dt2026/api/records"
	"dt2026/httpx"
	"dt2026/lib/optional"
	"dt2026/server/http"
)

const Name = "records_query"

const Description = "Search mdk's records by optional text phrases (q), tags (tag), and created_at range (from inclusive, to exclusive). Returns a page of matches. page defaults to 1; page_size defaults to 100 and must be 1–1000. At most 10 q values and 10 tags; empty q strings are rejected. Do not use this to create a record, fetch one by id, or as a health check."

type Input struct {
	Q        []string `json:"q,omitempty" jsonschema:"Optional search phrases; each matches raw_content, objective_context, ai_analysis, or tags (ILIKE). At most 10; an empty string is rejected. Omit for no text filter."`
	Tag      []string `json:"tag,omitempty" jsonschema:"Optional tags; all must be present (AND). At most 10; trimmed, no blanks, no duplicates. Omit for no tag filter."`
	From     *string  `json:"from,omitempty" jsonschema:"Optional inclusive lower bound on created_at. RFC 3339 or DateOnly (YYYY-MM-DD, Asia/Shanghai)."`
	To       *string  `json:"to,omitempty" jsonschema:"Optional exclusive upper bound on created_at. RFC 3339 or DateOnly (YYYY-MM-DD, Asia/Shanghai)."`
	Page     *int64   `json:"page,omitempty" jsonschema:"Optional page number, default 1, must be greater than 0."`
	PageSize *int64   `json:"page_size,omitempty" jsonschema:"Optional page size, default 100, must be between 1 and 1000."`
}

type Output = records.QueryRecordResponse

func RecordsQuery(s *http.Server, in Input) (Output, *httpx.Error) {
	criteria, err := in.criteria()
	if err != nil {
		var zero Output
		return zero, err
	}

	pool, err := s.DB()
	if err != nil {
		var zero Output
		return zero, err
	}

	res := records.Query(criteria, pool)
	if httpErr, ok := res.(*httpx.Error); ok {
		var zero Output
		return zero, httpErr
	}
	return res.(records.QueryRecordResponse), nil
}

func (in Input) criteria() (*records.QueryCriteria, *httpx.Error) {
	var from optional.Optional[string]
	if in.From != nil {
		from = optional.Some(*in.From)
	}
	parsedFrom, err := records.ParseOptionalJSONTime(from)
	if err != nil {
		return nil, err
	}

	var to optional.Optional[string]
	if in.To != nil {
		to = optional.Some(*in.To)
	}
	parsedTo, err := records.ParseOptionalJSONTime(to)
	if err != nil {
		return nil, err
	}

	var page optional.Optional[int64]
	if in.Page != nil {
		page = optional.Some(*in.Page)
	}
	var pageSize optional.Optional[int64]
	if in.PageSize != nil {
		pageSize = optional.Some(*in.PageSize)
	}

	return &records.QueryCriteria{
		Queries:  in.Q,
		Tags:     in.Tag,
		From:     parsedFrom,
		To:       parsedTo,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
