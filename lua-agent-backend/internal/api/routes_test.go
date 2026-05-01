package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"log/slog"

	"lua-agent/backend/internal/agent"
	"lua-agent/backend/internal/storage"
	"lua-agent/backend/internal/validator"
)

type stubAgent struct {
	output *agent.GenerateOutput
	err    error
}

func (s *stubAgent) Generate(_ context.Context, _ agent.GenerateInput) (*agent.GenerateOutput, error) {
	return s.output, s.err
}

func (s *stubAgent) PingLLM(_ context.Context) error {
	return nil
}

type stubRepo struct{}

func (s *stubRepo) Save(_ context.Context, _ storage.SaveHistoryInput) (*storage.DBHistory, error) {
	return nil, nil
}

func (s *stubRepo) GetRecentSuccess(_ context.Context, _ int) ([]storage.DBHistory, error) {
	return []storage.DBHistory{{ID: "1", UserPrompt: "demo", Success: true}}, nil
}

func (s *stubRepo) GetStats(_ context.Context, _ storage.StatsFilter) (*storage.Stats, error) {
	return &storage.Stats{TotalRuns: 1, SuccessfulRuns: 1, SuccessRate: 1}, nil
}

func (s *stubRepo) Ping(_ context.Context) error {
	return nil
}

func TestGenerateRoute(t *testing.T) {
	handler := NewRouter(&stubAgent{
		output: &agent.GenerateOutput{
			SessionID: "s1",
			Code:      "print(1)",
			Validation: validator.Result{
				OK: true,
			},
			Model: "demo-model",
		},
	}, &stubRepo{}, slog.Default())

	req := httptest.NewRequest(http.MethodPost, "/generate", strings.NewReader(`{"prompt":"demo lua task"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var response GenerateResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.Code != "print(1)" {
		t.Fatalf("unexpected code: %s", response.Code)
	}
}

func TestHealthRoute(t *testing.T) {
	handler := NewRouter(&stubAgent{output: &agent.GenerateOutput{}}, &stubRepo{}, slog.Default())

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), `"status":"ok"`) {
		t.Fatalf("unexpected health response: %s", rr.Body.String())
	}
}
