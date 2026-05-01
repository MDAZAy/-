package api

import "lua-agent/backend/internal/validator"

type GenerateRequest struct {
	SessionID string `json:"session_id"`
	Prompt    string `json:"prompt"`
}

type GenerateResponse struct {
	SessionID          string           `json:"session_id"`
	Code               string           `json:"code,omitempty"`
	Plan               string           `json:"plan,omitempty"`
	Validation         validator.Result `json:"validation"`
	NeedsClarification bool             `json:"needs_clarification"`
	Clarification      string           `json:"clarification,omitempty"`
	Corrected          bool             `json:"corrected"`
	Model              string           `json:"model,omitempty"`
}

type HistoryResponse struct {
	Items interface{} `json:"items"`
}

type StatsResponse struct {
	Stats interface{} `json:"stats"`
}
