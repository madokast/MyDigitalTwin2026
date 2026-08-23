package detach

import (
	"dt2026/api/records/tags"
	"dt2026/httpx"
	"dt2026/server/http"
)

const Name = "records_tags_detach"

const Description = "Detach one tag from one of mdk's records. record_id and tag are required; tag is trimmed and must not be blank. Returns the record's tags after the detach. Removing a present tag is detached true and changed true; repeating when it is already gone is detached false and changed false. Status is always 200 on success. Do not use this to attach a tag, create a record, or list every tag in the library."

type Input struct {
	RecordID int64  `json:"record_id" jsonschema:"Record ID assigned by the server when the record was created"`
	Tag      string `json:"tag" jsonschema:"Tag to detach; trimmed, must not be blank"`
}

type Output = tags.DetachTagResponse

func RecordsTagsDetach(s *http.Server, in Input) (Output, *httpx.Error) {
	pool, err := s.DB()
	if err != nil {
		var zero Output
		return zero, err
	}

	res := tags.Detach(in.RecordID, in.Tag, pool)
	if httpErr, ok := res.(*httpx.Error); ok {
		var zero Output
		return zero, httpErr
	}
	return res.(tags.DetachTagResponse), nil
}
