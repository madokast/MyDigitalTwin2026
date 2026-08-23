package get_all

import (
	"dt2026/api/records/tags"
	"dt2026/httpx"
	"dt2026/server/http"
)

const Name = "records_tags_get_all"

const Description = "List every tag that appears in mdk's records, with how many records carry each tag. Returns tag_counts sorted by count descending then tag ascending. No arguments. Unused tags are omitted; if none exist, tag_counts is []. Do not use this to search records, fetch a record by id, or as a health check."

type Input struct {
}

type Output = tags.GetAllTagsResponse

func RecordsTagsGetAll(s *http.Server, _ Input) (Output, *httpx.Error) {
	pool, err := s.DB()
	if err != nil {
		var zero Output
		return zero, err
	}

	res := tags.GetAll(pool)
	if httpErr, ok := res.(*httpx.Error); ok {
		var zero Output
		return zero, httpErr
	}
	return res.(tags.GetAllTagsResponse), nil
}
