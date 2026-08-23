package records

import (
	"context"
	"dt2026/api/notify/qqbot"
	"dt2026/api/records/tags"
	"dt2026/httpx"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostResponse struct {
	Ok     httpx.TrueType `json:"ok" jsonschema:"Whether the request succeeded"`
	Status int    `json:"status" jsonschema:"HTTP status code"`
	Record Record `json:"record" jsonschema:"The created record"`
}

func (r PostResponse) GetStatus() int {
	return r.Status
}

const postRecordNotifyMessage = `New record created: %d
Raw content: %s
Objective context: %s
Ai analysis: %s
Tags: %s`

func Post(record *NewRecord, pool *pgxpool.Pool, bot *qqbot.Sender) httpx.Response {
	if err := normalizeNewRecord(record); err != nil {
		return err
	}

	var response = PostResponse{
		Ok:     httpx.True,
		Status: http.StatusCreated,
		Record: Record{
			CreatedAt:        JSONTime(time.Now().UTC()),
			RawContent:       record.RawContent,
			ObjectiveContext: record.ObjectiveContext,
			AIAnalysis:       record.AIAnalysis,
			Tags:             record.Tags,
		},
	}

	err := pool.QueryRow(context.Background(), insertRecordSQL,
		response.Record.CreatedAt.GoTime(),
		response.Record.RawContent,
		response.Record.ObjectiveContext,
		response.Record.AIAnalysis,
		response.Record.Tags,
	).Scan(&response.Record.ID)

	if err != nil {
		return httpx.NewInternalServerError(fmt.Sprintf(
			"failed to insert record while executing %s with %v: %s",
			insertRecordSQL,
			[]any{
				response.Record.CreatedAt.GoTime(),
				response.Record.RawContent,
				response.Record.ObjectiveContext,
				response.Record.AIAnalysis,
				response.Record.Tags,
			},
			err.Error(),
		))
	}

	bot.SendMessageAsync(fmt.Sprintf(postRecordNotifyMessage,
		response.Record.ID,
		response.Record.RawContent,
		response.Record.ObjectiveContext,
		response.Record.AIAnalysis,
		"["+strings.Join(response.Record.Tags, ", ")+"]",
	))

	return response
}

func normalizeNewRecord(record *NewRecord) *httpx.Error {
	if strings.TrimSpace(record.RawContent) == "" {
		return httpx.NewBadRequestError("raw content cannot be empty or contain only whitespace")
	}

	if strings.TrimSpace(record.ObjectiveContext) == "" {
		return httpx.NewBadRequestError("objective context cannot be empty or contain only whitespace")
	}

	if strings.TrimSpace(record.AIAnalysis) == "" {
		return httpx.NewBadRequestError("AI analysis cannot be empty or contain only whitespace")
	}

	var err *httpx.Error
	record.Tags, err = tags.NormalizeTags(record.Tags)
	if err != nil {
		return err
	}

	// 空 tags 不能是 nil，必须 []
	// null value in column \"tags\" of relation \"records\" violates not-null constraint (SQLSTATE 23502)
	if len(record.Tags) == 0 {
		record.Tags = []string{}
	}

	return nil
}
