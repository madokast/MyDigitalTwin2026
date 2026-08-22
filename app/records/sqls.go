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
SELECT id, created_at, raw_content, 
  objective_context, ai_analysis, 
  tags, COUNT(*) OVER() AS total
FROM records WHERE 1=1
`

const exportRecordSQL = `
SELECT id, created_at, raw_content, 
  objective_context, ai_analysis, 
  tags
FROM records WHERE 1=1
`

const countRecordSQL = `
SELECT COUNT(*) AS total
FROM records WHERE 1=1
`

var textColumns = []string{
	"raw_content",
	"objective_context",
	"ai_analysis",
}
