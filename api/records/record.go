package records

type Record struct {
	ID               int64    `json:"id" jsonschema:"Record ID assigned by the server"`
	CreatedAt        JSONTime `json:"created_at" jsonschema:"Creation time in RFC 3339 with millisecond precision and a +08:00 offset"`
	RawContent       string   `json:"raw_content" jsonschema:"mdk's raw utterance, recorded as-is"`
	ObjectiveContext string   `json:"objective_context" jsonschema:"Objective conversation context"`
	AIAnalysis       string   `json:"ai_analysis" jsonschema:"AI analysis"`
	Tags             []string `json:"tags" jsonschema:"Tags on this record"`
}

type NewRecord struct {
	RawContent       string   `json:"raw_content" jsonschema:"mdk's raw utterance; required, must not be blank"`
	ObjectiveContext string   `json:"objective_context" jsonschema:"Objective conversation context; required, must not be blank"`
	AIAnalysis       string   `json:"ai_analysis" jsonschema:"AI analysis; required, must not be blank"`
	Tags             []string `json:"tags,omitempty" jsonschema:"Optional tags; trimmed, no blanks, no duplicates. Omit or pass an empty array for none."`
}
