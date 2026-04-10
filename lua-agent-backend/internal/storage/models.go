package storage

import "time"

const (
	ValidationStatusUnknown = "unknown"
	ValidationStatusPassed  = "passed"
	ValidationStatusFailed  = "failed"
)

// DBSession mirrors the sessions table and keeps conversation continuity.
type DBSession struct {
	ID         string
	LastPrompt string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// DBHistory stores one generation attempt and its validation result.
type DBHistory struct {
	ID               string
	SessionID        string
	UserPrompt       string
	ClarifiedPrompt  string
	GeneratedCode    string
	ValidationStatus string
	ValidationErrors []string
	Success          bool
	ModelName        string
	LatencyMS        int64
	InputTokens      int
	OutputTokens     int
	Metadata         map[string]any
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// SaveHistoryInput is the payload accepted by Repository.Save.
type SaveHistoryInput struct {
	SessionID        string
	UserPrompt       string
	ClarifiedPrompt  string
	GeneratedCode    string
	ValidationStatus string
	ValidationErrors []string
	Success          bool
	ModelName        string
	LatencyMS        int64
	InputTokens      int
	OutputTokens     int
	Metadata         map[string]any
}

// Stats contains jury-friendly aggregate metrics for the generation pipeline.
type Stats struct {
	TotalRuns          int64
	SuccessfulRuns     int64
	FailedRuns         int64
	SuccessRate        float64
	AverageLatencyMS   float64
	ValidationStatuses map[string]int64
}
