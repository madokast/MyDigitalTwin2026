package records

import "time"

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

const insertRecordNotifyMessage = `New record created: %d
Raw content: %s
Objective context: %s
Ai analysis: %s
Tags: %s`

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

func (t *JSONTime) Time() time.Time {
	return time.Time(*t)
}

var loc = time.FixedZone("CST", 8*3600)

func (t JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(
		`"` + time.Time(t).In(loc).Format(time.RFC3339) + `"`,
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
