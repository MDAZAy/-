package agent

import (
	"context"
	"strings"

	"lua-agent/backend/internal/llm"
	"lua-agent/backend/internal/validator"
)

type Corrector struct {
	client *llm.Client
	config LLMConfig
}

func NewCorrector(client *llm.Client, cfg LLMConfig) *Corrector {
	return &Corrector{
		client: client,
		config: cfg,
	}
}

func (c *Corrector) Correct(ctx context.Context, prompt string, code string, validationResult validator.Result) string {
	if c == nil || c.client == nil {
		return code
	}

	resp, err := c.client.Generate(ctx, llm.GenerateInput{
		Model:       c.config.Model,
		System:      llm.CorrectionPrompt,
		Prompt:      strings.TrimSpace("Original task:\n" + prompt + "\n\nOriginal code:\n" + code + "\n\nValidation issues:\n" + joinIssues(validationResult.Issues) + "\n\nReturn corrected Lua code only."),
		NumCtx:      c.config.NumCtx,
		NumPredict:  c.config.NumPredict,
		Batch:       c.config.Batch,
		Parallel:    c.config.Parallel,
		Temperature: 0.2,
	})
	if err != nil || strings.TrimSpace(resp.Text) == "" {
		return code
	}

	return strings.TrimSpace(resp.Text)
}
