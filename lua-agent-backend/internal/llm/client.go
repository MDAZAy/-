package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	endpoint   string
	model      string
	httpClient *http.Client
}

type GenerateInput struct {
	Model       string
	Prompt      string
	System      string
	NumCtx      int
	NumPredict  int
	Batch       int
	Parallel    int
	Temperature float64
}

type GenerateOutput struct {
	Text              string
	Model             string
	PromptEvalCount   int
	EvalCount         int
	TotalDuration     int64
	LoadDuration      int64
	PromptEvalNanos   int64
	EvalDurationNanos int64
}

type ollamaGenerateRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	System  string                 `json:"system,omitempty"`
	Stream  bool                   `json:"stream"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type ollamaGenerateResponse struct {
	Model             string `json:"model"`
	Response          string `json:"response"`
	PromptEvalCount   int    `json:"prompt_eval_count"`
	EvalCount         int    `json:"eval_count"`
	TotalDuration     int64  `json:"total_duration"`
	LoadDuration      int64  `json:"load_duration"`
	PromptEvalNanos   int64  `json:"prompt_eval_duration"`
	EvalDurationNanos int64  `json:"eval_duration"`
}

func NewClient(endpoint string, model string, timeout time.Duration) *Client {
	if timeout <= 0 {
		timeout = 60 * time.Second
	}

	return &Client{
		endpoint:   strings.TrimRight(endpoint, "/"),
		model:      model,
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (c *Client) Generate(ctx context.Context, input GenerateInput) (*GenerateOutput, error) {
	model := strings.TrimSpace(input.Model)
	if model == "" {
		model = c.model
	}
	if model == "" {
		return nil, fmt.Errorf("llm model is required")
	}
	if strings.TrimSpace(input.Prompt) == "" {
		return nil, fmt.Errorf("llm prompt is required")
	}

	reqBody := ollamaGenerateRequest{
		Model:  model,
		Prompt: input.Prompt,
		System: input.System,
		Stream: false,
		Options: map[string]interface{}{
			"num_ctx":     fallbackInt(input.NumCtx, 4096),
			"num_predict": fallbackInt(input.NumPredict, 256),
			"batch":       fallbackInt(input.Batch, 1),
			"parallel":    fallbackInt(input.Parallel, 1),
			"temperature": input.Temperature,
		},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal ollama request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint+"/api/generate", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build ollama request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call ollama generate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("ollama generate returned status %s", resp.Status)
	}

	var decoded ollamaGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("decode ollama response: %w", err)
	}

	return &GenerateOutput{
		Text:              strings.TrimSpace(decoded.Response),
		Model:             decoded.Model,
		PromptEvalCount:   decoded.PromptEvalCount,
		EvalCount:         decoded.EvalCount,
		TotalDuration:     decoded.TotalDuration,
		LoadDuration:      decoded.LoadDuration,
		PromptEvalNanos:   decoded.PromptEvalNanos,
		EvalDurationNanos: decoded.EvalDurationNanos,
	}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint+"/api/tags", nil)
	if err != nil {
		return fmt.Errorf("build ollama ping request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("call ollama ping: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("ollama ping returned status %s", resp.Status)
	}

	return nil
}

func fallbackInt(value int, fallback int) int {
	if value > 0 {
		return value
	}

	return fallback
}
