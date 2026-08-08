package records

const CreateRecordSQL = `
CREATE TABLE IF NOT EXISTS records (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,

    raw_content TEXT NOT NULL,
    objective_context TEXT NOT NULL,
    ai_analysis TEXT NOT NULL,

    tags TEXT[] NOT NULL
);
`

const insertRecordSQL = `
INSERT INTO records
(created_at, raw_content, objective_context, ai_analysis, tags)
VALUES ($1,$2,$3,$4,$5)
RETURNING id;
`

const queryRecordSQL = `
SELECT id, created_at, raw_content, objective_context, ai_analysis, tags
FROM records WHERE 1=1
`

const insertRecordNotifyMessage = `New record created: %d
Raw content: %s
Objective context: %s
Ai analysis: %s
Tags: %s`

type Record struct {
	ID               int64    `json:"id"`
	CreatedAt        JSONTime `json:"created_at"`
	RawContent       string   `json:"raw_content"`
	ObjectiveContext string   `json:"objective_context"`
	AIAnalysis       string   `json:"ai_analysis,omitempty"`
	Tags             []string `json:"tags"`
}

type NewRecord struct {
	RawContent       string   `json:"raw_content"`
	ObjectiveContext string   `json:"objective_context"`
	AIAnalysis       string   `json:"ai_analysis,omitempty"`
	Tags             []string `json:"tags"`
}
