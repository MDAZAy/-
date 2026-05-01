package agent

import (
	"context"
	"strings"
	"time"

	"lua-agent/backend/internal/llm"
	"lua-agent/backend/internal/storage"
	"lua-agent/backend/internal/validator"
)

const ClarifyPrompt = `Ask one short clarification question when the task is too vague to generate correct Lua code.`

type Generator interface {
	Generate(ctx context.Context, input llm.GenerateInput) (*llm.GenerateOutput, error)
}

type Repository interface {
	Save(ctx context.Context, input storage.SaveHistoryInput) (*storage.DBHistory, error)
	GetRecentSuccess(ctx context.Context, limit int) ([]storage.DBHistory, error)
	GetStats(ctx context.Context, filter storage.StatsFilter) (*storage.Stats, error)
}

type LLMConfig struct {
	Model      string
	NumCtx     int
	NumPredict int
	Batch      int
	Parallel   int
}

type Service struct {
	generator Generator
	repo      Repository
	validate  *validator.Validator
	planner   *Planner
	corrector *Corrector
	config    LLMConfig
}

type GenerateInput struct {
	SessionID string
	Prompt    string
}

type GenerateOutput struct {
	SessionID          string           `json:"session_id"`
	Code               string           `json:"code,omitempty"`
	Plan               string           `json:"plan,omitempty"`
	Validation         validator.Result `json:"validation"`
	NeedsClarification bool             `json:"needs_clarification"`
	Clarification      string           `json:"clarification,omitempty"`
	Corrected          bool             `json:"corrected"`
	Model              string           `json:"model,omitempty"`
}

func NewService(generator Generator, repo Repository, validate *validator.Validator, cfg LLMConfig) *Service {
	return &Service{
		generator: generator,
		repo:      repo,
		validate:  validate,
		planner:   NewPlanner(asLLMClient(generator), cfg),
		corrector: NewCorrector(asLLMClient(generator), cfg),
		config:    cfg,
	}
}

func (s *Service) Generate(ctx context.Context, input GenerateInput) (*GenerateOutput, error) {
	prompt := strings.TrimSpace(input.Prompt)
	if needsClarification(prompt) {
		return &GenerateOutput{
			SessionID:          input.SessionID,
			NeedsClarification: true,
			Clarification:      buildClarificationQuestion(prompt),
			Validation:         validator.Result{OK: false},
		}, nil
	}

	plan := s.planner.Plan(ctx, prompt)
	started := time.Now()

	generated, err := s.generator.Generate(ctx, llm.GenerateInput{
		Model:       s.config.Model,
		System:      llm.SystemPrompt,
		Prompt:      strings.TrimSpace("Task:\n" + prompt + "\n\nPlan:\n" + plan + "\n\nReturn only Lua code without markdown fences."),
		NumCtx:      s.config.NumCtx,
		NumPredict:  s.config.NumPredict,
		Batch:       s.config.Batch,
		Parallel:    s.config.Parallel,
		Temperature: 0.2,
	})
	if err != nil {
		return nil, err
	}

	code := stripMarkdownCodeFence(generated.Text)
	validationResult := s.validate.Validate(ctx, code)
	corrected := false

	if !validationResult.OK {
		correctedCode := stripMarkdownCodeFence(s.corrector.Correct(ctx, prompt, code, validationResult))
		if correctedCode != "" && correctedCode != code {
			corrected = true
			code = correctedCode
			validationResult = s.validate.Validate(ctx, code)
		}
	}

	if s.repo != nil {
		_, _ = s.repo.Save(ctx, storage.SaveHistoryInput{
			SessionID:        input.SessionID,
			UserPrompt:       prompt,
			ClarifiedPrompt:  prompt,
			GeneratedCode:    code,
			ValidationStatus: validationStatus(validationResult.OK),
			ValidationErrors: extractValidationErrors(validationResult),
			Success:          validationResult.OK,
			ModelName:        generated.Model,
			LatencyMS:        time.Since(started).Milliseconds(),
			InputTokens:      generated.PromptEvalCount,
			OutputTokens:     generated.EvalCount,
			Metadata: map[string]any{
				"plan":       plan,
				"corrected":  corrected,
				"sandbox_ok": validationResult.Sandbox != nil && validationResult.Sandbox.OK,
			},
		})
	}

	return &GenerateOutput{
		SessionID:  input.SessionID,
		Code:       code,
		Plan:       plan,
		Validation: validationResult,
		Corrected:  corrected,
		Model:      generated.Model,
	}, nil
}

func (s *Service) PingLLM(ctx context.Context) error {
	pinger, ok := s.generator.(interface {
		Ping(context.Context) error
	})
	if !ok {
		return nil
	}

	return pinger.Ping(ctx)
}

func buildClarificationQuestion(prompt string) string {
	if prompt == "" {
		return "Что должен делать Lua-скрипт: какие входы, какой ожидаемый результат и в какой среде он будет запускаться?"
	}

	return "Уточни, пожалуйста, входные данные, ожидаемый результат и среду выполнения Lua-скрипта, чтобы я сгенерировал корректный код."
}

func needsClarification(prompt string) bool {
	if len([]rune(prompt)) < 12 {
		return true
	}

	low := strings.ToLower(prompt)
	keywords := []string{"lua", "скрипт", "script", "function", "таблиц", "json", "mws", "wf.", "_utils"}
	for _, keyword := range keywords {
		if strings.Contains(low, keyword) {
			return false
		}
	}

	return len(strings.Fields(prompt)) < 4
}

func stripMarkdownCodeFence(code string) string {
	code = strings.TrimSpace(code)
	code = strings.TrimPrefix(code, "```lua")
	code = strings.TrimPrefix(code, "```")
	code = strings.TrimSuffix(code, "```")
	return strings.TrimSpace(code)
}

func validationStatus(ok bool) string {
	if ok {
		return storage.ValidationStatusPassed
	}

	return storage.ValidationStatusFailed
}

func extractValidationErrors(result validator.Result) []string {
	errors := make([]string, 0, len(result.Issues))
	for _, issue := range result.Issues {
		errors = append(errors, issue.Kind+": "+issue.Message)
	}

	return errors
}

func joinIssues(issues []validator.Issue) string {
	if len(issues) == 0 {
		return "no validation issues reported"
	}

	parts := make([]string, 0, len(issues))
	for _, issue := range issues {
		parts = append(parts, string(issue.Level)+" "+issue.Kind+": "+issue.Message)
	}

	return strings.Join(parts, "\n")
}

func asLLMClient(generator Generator) *llm.Client {
	client, ok := generator.(*llm.Client)
	if !ok {
		return nil
	}

	return client
}

func minInt(a int, b int) int {
	if a == 0 || a > b {
		return b
	}

	return a
}
