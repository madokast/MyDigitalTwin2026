package records

import "time"

const createRecordSQL = `
CREATE TABLE records (
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

type JSONTime time.Time

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

func (t JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(
		`"` + time.Time(t).Format(time.RFC3339) + `"`,
	), nil
}

func (t *JSONTime) UnmarshalJSON(data []byte) error {
	parsed, err := time.Parse(
		time.RFC3339,
		string(data),
	)
	if err != nil {
		return err
	}

	*t = JSONTime(parsed)
	return nil
}
