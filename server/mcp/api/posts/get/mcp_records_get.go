package get

import (
	"dt2026/api/records"
	"dt2026/httpx"
	"dt2026/server/http"
)

const Name = "records_get"

const Description = "Fetch one of mdk's records by id. Returns the raw utterance, objective context, AI analysis, tags, and creation time. record_id is required. Do not use this to create records, list history, or as a health check."

type Input struct {
	RecordID int64 `json:"record_id" jsonschema:"Record ID assigned by the server when the record was created"`
}

type Output = records.GetRecordsResponse

func RecordsGet(s *http.Server, in Input) (Output, *httpx.Error) {
	pool, err := s.DB()
	if err != nil {
		var zero Output
		return zero, err
	}

	res := records.Get(in.RecordID, pool)
	if httpErr, ok := res.(*httpx.Error); ok {
		var zero Output
		return zero, httpErr
	}
	return res.(records.GetRecordsResponse), nil
}
