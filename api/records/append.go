package records

import (
	"context"
	"dt2026/api/notify/qqbot"
	"dt2026/httpx"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AppendRecordResponse struct {
	Ok     bool   `json:"ok"`
	Status int    `json:"status"`
	Record Record `json:"record"`
}

func (r AppendRecordResponse) GetStatus() int {
	return r.Status
}

func Append(record *NewRecord, pool *pgxpool.Pool, bot *qqbot.Sender) httpx.Response {
	if err := normalizeNewRecord(record); err != nil {
		return err
	}

	var response = AppendRecordResponse{
		Ok:     true,
		Status: http.StatusOK,
		Record: Record{
			CreatedAt:        JSONTime(time.Now().UTC()),
			RawContent:       record.RawContent,
			ObjectiveContext: record.ObjectiveContext,
			AIAnalysis:       record.AIAnalysis,
			Tags:             record.Tags,
		},
	}

	err := pool.QueryRow(context.Background(), insertRecordSQL,
		response.Record.CreatedAt.Time(),
		response.Record.RawContent,
		response.Record.ObjectiveContext,
		response.Record.AIAnalysis,
		response.Record.Tags,
	).Scan(&response.Record.ID)

	if err != nil {
		return httpx.NewInternalServerError("failed to insert record: " + err.Error())
	}

	go func() {
		_ = bot.SendMessage(fmt.Sprintf(insertRecordNotifyMessage,
			response.Record.ID,
			response.Record.RawContent,
			response.Record.ObjectiveContext,
			response.Record.AIAnalysis,
			"["+strings.Join(response.Record.Tags, ", ")+"]",
		))
	}()

	return response
}

func normalizeNewRecord(record *NewRecord) *httpx.Error {
	record.RawContent = strings.TrimSpace(record.RawContent)
	if record.RawContent == "" {
		return httpx.NewBadRequestError("raw content cannot be empty")
	}

	record.ObjectiveContext = strings.TrimSpace(record.ObjectiveContext)
	if record.ObjectiveContext == "" {
		return httpx.NewBadRequestError("objective context cannot be empty")
	}

	record.AIAnalysis = strings.TrimSpace(record.AIAnalysis)
	if record.AIAnalysis == "" {
		return httpx.NewBadRequestError("AI analysis cannot be empty")
	}

	var tags []string
	var seenTags = make(map[string]bool)
	for _, tag := range record.Tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			return httpx.NewBadRequestError("tags cannot contain empty strings")
		}
		if seenTags[tag] {
			return httpx.NewBadRequestError("duplicate tags are not allowed")
		}
		seenTags[tag] = true
		tags = append(tags, tag)
	}
	// null value in column \"tags\" of relation \"records\" violates not-null constraint (SQLSTATE 23502)
	if len(tags) == 0 {
		tags = []string{}
	}
	record.Tags = tags

	return nil
}
