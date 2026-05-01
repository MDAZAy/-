package agent

import (
	"context"
	"strings"

	"lua-agent/backend/internal/llm"
)

type Planner struct {
	client *llm.Client
	config LLMConfig
}

func NewPlanner(client *llm.Client, cfg LLMConfig) *Planner {
	return &Planner{
		client: client,
		config: cfg,
	}
}

func (p *Planner) Plan(ctx context.Context, prompt string) string {
	if p == nil || p.client == nil {
		return fallbackPlan(prompt)
	}

	resp, err := p.client.Generate(ctx, llm.GenerateInput{
		Model:       p.config.Model,
		Prompt:      "User task:\n" + prompt,
		System:      llm.CotPrompt,
		NumCtx:      p.config.NumCtx,
		NumPredict:  minInt(p.config.NumPredict, 192),
		Batch:       p.config.Batch,
		Parallel:    p.config.Parallel,
		Temperature: 0.1,
	})
	if err != nil || strings.TrimSpace(resp.Text) == "" {
		return fallbackPlan(prompt)
	}

	return strings.TrimSpace(resp.Text)
}

func fallbackPlan(prompt string) string {
	return "1. Parse the requested behavior.\n2. Generate a safe Lua solution.\n3. Validate syntax, security, and runtime behavior.\n4. Correct the script if validation fails.\nTask: " + strings.TrimSpace(prompt)
}
