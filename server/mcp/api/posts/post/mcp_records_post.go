package post

import (
	"dt2026/api/records"
	"dt2026/httpx"
	"dt2026/server/http"
)

const Name = "records_post"

const Description = "Create a record of mdk's utterance with objective context, AI analysis, and optional tags. raw_content, objective_context, and ai_analysis are required and must not be blank. tags are optional; omit or pass [] for none. Do not use this as a health check or to query history."

type Input = records.NewRecord

type Output = records.PostResponse

func RecordsPost(s *http.Server, in Input) (Output, *httpx.Error) {
	pool, err := s.DB()
	if err != nil {
		var zero Output
		return zero, err
	}
	bot, err := s.QQBot()
	if err != nil {
		var zero Output
		return zero, err
	}

	res := records.Post(&in, pool, bot)
	if httpErr, ok := res.(*httpx.Error); ok {
		var zero Output
		return zero, httpErr
	}
	return res.(records.PostResponse), nil
}
