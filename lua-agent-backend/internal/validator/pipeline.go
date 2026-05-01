package validator

import (
	"context"
	"strings"
	"time"
)

type Level string

const (
	LevelError   Level = "error"
	LevelWarning Level = "warning"
)

type Issue struct {
	Kind    string `json:"kind"`
	Level   Level  `json:"level"`
	Message string `json:"message"`
}

type SandboxResult struct {
	OK       bool          `json:"ok"`
	Stdout   string        `json:"stdout"`
	Error    string        `json:"error,omitempty"`
	TimedOut bool          `json:"timed_out"`
	Duration time.Duration `json:"duration"`
}

type Result struct {
	OK      bool           `json:"ok"`
	Issues  []Issue        `json:"issues"`
	Sandbox *SandboxResult `json:"sandbox,omitempty"`
}

type Validator struct {
	syntax   *SyntaxValidator
	security *SecurityValidator
	mws      *MWSValidator
	sandbox  *SandboxValidator
}

func New(timeout time.Duration) *Validator {
	return &Validator{
		syntax:   NewSyntaxValidator(),
		security: NewSecurityValidator(),
		mws:      NewMWSValidator(),
		sandbox:  NewSandboxValidator(timeout),
	}
}

func (v *Validator) Validate(ctx context.Context, code string) Result {
	result := Result{
		OK:     true,
		Issues: make([]Issue, 0, 8),
	}

	code = strings.TrimSpace(code)
	if code == "" {
		result.OK = false
		result.Issues = append(result.Issues, Issue{
			Kind:    "input",
			Level:   LevelError,
			Message: "empty lua code",
		})
		return result
	}

	result.Issues = append(result.Issues, v.syntax.Validate(code)...)
	result.Issues = append(result.Issues, v.security.Validate(code)...)
	result.Issues = append(result.Issues, v.mws.Validate(code)...)

	if hasBlockingIssues(result.Issues) {
		result.OK = false
		return result
	}

	sandbox := v.sandbox.Validate(ctx, code)
	result.Sandbox = sandbox
	if !sandbox.OK {
		result.OK = false
		result.Issues = append(result.Issues, Issue{
			Kind:    "sandbox",
			Level:   LevelError,
			Message: sandbox.Error,
		})
	}

	return result
}

func hasBlockingIssues(issues []Issue) bool {
	for _, issue := range issues {
		if issue.Level == LevelError {
			return true
		}
	}

	return false
}
