package attach

import (
	"dt2026/api/records/tags"
	"dt2026/httpx"
	"dt2026/server/http"
)

const Name = "records_tags_attach"

const Description = "Attach one tag to one of mdk's records. record_id and tag are required; tag is trimmed and must not be blank. Returns the record's tags after the attach. A new tag is status 201 and changed true; repeating the same tag is 200 and changed false. Do not use this to create a record, fetch a record by id, or list every tag in the library."

type Input struct {
	RecordID int64  `json:"record_id" jsonschema:"Record ID assigned by the server when the record was created"`
	Tag      string `json:"tag" jsonschema:"Tag to attach; trimmed, must not be blank"`
}

type Output = tags.AttachTagResponse

func RecordsTagsAttach(s *http.Server, in Input) (Output, *httpx.Error) {
	pool, err := s.DB()
	if err != nil {
		var zero Output
		return zero, err
	}

	res := tags.Attach(in.RecordID, in.Tag, pool)
	if httpErr, ok := res.(*httpx.Error); ok {
		var zero Output
		return zero, httpErr
	}
	return res.(tags.AttachTagResponse), nil
}
