package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"lua-agent/backend/internal/agent"
	"lua-agent/backend/internal/storage"
)

type agentService interface {
	Generate(ctx context.Context, input agent.GenerateInput) (*agent.GenerateOutput, error)
	PingLLM(ctx context.Context) error
}

type Router struct {
	agent  agentService
	repo   storage.Repository
	logger *slog.Logger
}

func NewRouter(agentService agentService, repo storage.Repository, logger *slog.Logger) http.Handler {
	router := &Router{
		agent:  agentService,
		repo:   repo,
		logger: logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", router.handleHealth)
	mux.HandleFunc("POST /generate", router.handleGenerate)
	mux.HandleFunc("GET /history", router.handleHistory)
	mux.HandleFunc("GET /stats", router.handleStats)
	if _, err := os.Stat("web/index.html"); err == nil {
		mux.Handle("/", http.FileServer(http.Dir("web")))
	}
	return router.withLogging(mux)
}

func (r *Router) handleHealth(w http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	components := map[string]map[string]string{
		"storage": {"status": "ok"},
		"ollama":  {"status": "ok"},
	}
	statusCode := http.StatusOK
	overall := "ok"

	if err := r.repo.Ping(ctx); err != nil {
		components["storage"] = map[string]string{
			"status": "error",
			"error":  err.Error(),
		}
		overall = "degraded"
		statusCode = http.StatusServiceUnavailable
	}

	if err := r.agent.PingLLM(ctx); err != nil {
		components["ollama"] = map[string]string{
			"status": "error",
			"error":  err.Error(),
		}
		overall = "degraded"
		statusCode = http.StatusServiceUnavailable
	}

	writeJSON(w, statusCode, map[string]any{
		"status":     overall,
		"components": components,
	})
}

func (r *Router) handleGenerate(w http.ResponseWriter, req *http.Request) {
	var body GenerateRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	result, err := r.agent.Generate(req.Context(), agent.GenerateInput{
		SessionID: body.SessionID,
		Prompt:    body.Prompt,
	})
	if err != nil {
		r.logger.Error("generate failed", "error", err)
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, GenerateResponse{
		SessionID:          result.SessionID,
		Code:               result.Code,
		Plan:               result.Plan,
		Validation:         result.Validation,
		NeedsClarification: result.NeedsClarification,
		Clarification:      result.Clarification,
		Corrected:          result.Corrected,
		Model:              result.Model,
	})
}

func (r *Router) handleHistory(w http.ResponseWriter, req *http.Request) {
	limit := parseInt(req.URL.Query().Get("limit"), 10)
	items, err := r.repo.GetRecentSuccess(req.Context(), limit)
	if err != nil {
		r.logger.Error("history failed", "error", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, HistoryResponse{Items: items})
}

func (r *Router) handleStats(w http.ResponseWriter, req *http.Request) {
	filter := storage.StatsFilter{
		From: parseTimePtr(req.URL.Query().Get("from")),
		To:   parseTimePtr(req.URL.Query().Get("to")),
	}

	stats, err := r.repo.GetStats(req.Context(), filter)
	if err != nil {
		r.logger.Error("stats failed", "error", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, StatsResponse{Stats: stats})
}

func (r *Router) withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, req)
		r.logger.Info("request completed", "method", req.Method, "path", req.URL.Path, "duration_ms", time.Since(started).Milliseconds())
	})
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func parseInt(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return fallback
	}

	return value
}

func parseTimePtr(raw string) *time.Time {
	if raw == "" {
		return nil
	}

	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil
	}

	return &t
}
