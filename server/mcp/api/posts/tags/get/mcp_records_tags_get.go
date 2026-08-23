package get

import (
	"dt2026/api/records/tags"
	"dt2026/httpx"
	"dt2026/server/http"
)

const Name = "records_tags_get"

const Description = "Fetch the tags on one of mdk's records by id. Returns tags as a string array; empty if the record has none. record_id is required. Do not use this to list every tag in the library, search records, or as a health check."

type Input struct {
	RecordID int64 `json:"record_id" jsonschema:"Record ID assigned by the server when the record was created"`
}

type Output = tags.GetTagResponse

func RecordsTagsGet(s *http.Server, in Input) (Output, *httpx.Error) {
	pool, err := s.DB()
	if err != nil {
		var zero Output
		return zero, err
	}

	res := tags.Get(in.RecordID, pool)
	if httpErr, ok := res.(*httpx.Error); ok {
		var zero Output
		return zero, httpErr
	}
	return res.(tags.GetTagResponse), nil
}
