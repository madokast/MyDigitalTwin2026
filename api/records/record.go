package records

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
