package rename

import (
	"dt2026/api/records/tags"
	"dt2026/httpx"
	"dt2026/server/http"
)

const Name = "records_tags_rename"

const Description = "Rename one tag on one of mdk's records. record_id, tag, and new_tag are required; both tags are trimmed and must not be blank. Returns the record's tags after the rename. A real rename is status 200 and changed true; renaming to the same name is changed false; renaming an absent tag is renamed false and changed false; renaming onto an existing tag on the same record fails with status 409. Do not use this to attach, detach, or list tags."

type Input struct {
	RecordID int64  `json:"record_id" jsonschema:"Record ID assigned by the server when the record was created"`
	Tag      string `json:"tag" jsonschema:"Tag to rename; trimmed, must not be blank"`
	NewTag   string `json:"new_tag" jsonschema:"New name for the tag; trimmed, must not be blank"`
}

type Output = tags.RenameTagResponse

func RecordsTagsRename(s *http.Server, in Input) (Output, *httpx.Error) {
	pool, err := s.DB()
	if err != nil {
		var zero Output
		return zero, err
	}

	res := tags.Rename(in.RecordID, in.Tag, in.NewTag, pool)
	if httpErr, ok := res.(*httpx.Error); ok {
		var zero Output
		return zero, httpErr
	}
	return res.(tags.RenameTagResponse), nil
}
