package records

import (
	"dt2026/api/envkeys"
	"dt2026/api/httpx"
	"dt2026/env"
	"encoding/json/v2"
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

func (r AppendRecordResponse) IsOk() bool {
	return r.Ok
}

func (r AppendRecordResponse) GetStatus() int {
	return r.Status
}

func Append(w http.ResponseWriter, r *http.Request) {
	var record NewRecord

	// 获取和校验 NewRecord
	err := json.UnmarshalRead(r.Body, &record)
	if err != nil {
		httpx.WriteJSON(w, httpx.NewInternalServerError("failed to decode request body: "+err.Error()))
		return
	}

	if err := normalizeNewRecord(&record); err != nil {
		httpx.WriteJSON(w, err)
		return
	}

	url, ok := env.Get(envkeys.DatabaseUrl)
	if !ok {
		httpx.WriteJSON(w, httpx.NewInternalServerError("database url not set"))
		return
	}

	pool, err := pgxpool.New(r.Context(), url)
	if err != nil {
		httpx.WriteJSON(w, httpx.NewInternalServerError("failed to create database pool: "+err.Error()))
		return
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

	_, err = pool.Exec(r.Context(), createRecordSQL)
	if err != nil {
		httpx.WriteJSON(w, httpx.NewInternalServerError("failed to create records table: "+err.Error()))
		return
	}

	err = pool.QueryRow(r.Context(), insertRecordSQL,
		response.Record.CreatedAt.Time(),
		response.Record.RawContent,
		response.Record.ObjectiveContext,
		response.Record.AIAnalysis,
		response.Record.Tags,
	).Scan(&response.Record.ID)

	if err != nil {
		httpx.WriteJSON(w, httpx.NewInternalServerError("failed to insert record: "+err.Error()))
		return
	}

	httpx.WriteJSON(w, response)
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
