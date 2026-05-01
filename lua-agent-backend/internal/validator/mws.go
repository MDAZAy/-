package validator

import (
	"regexp"
	"strings"
)

var (
	wfVarsAccessPattern   = regexp.MustCompile(`wf\.vars\s*\[`)
	wfInitCallPattern     = regexp.MustCompile(`wf\.initVariables\s*\(`)
	utilsArrayCallPattern = regexp.MustCompile(`_utils\.array\s*\(`)
)

// MWSValidator applies lightweight domain checks for workflow-style Lua snippets.
type MWSValidator struct{}

func NewMWSValidator() *MWSValidator {
	return &MWSValidator{}
}

func (v *MWSValidator) Validate(code string) []Issue {
	var issues []Issue

	hasWFVars := wfVarsAccessPattern.MatchString(code)
	hasWFInit := wfInitCallPattern.MatchString(code)
	hasWFNamespace := strings.Contains(code, "wf.")
	hasUtilsNamespace := strings.Contains(code, "_utils.")

	if hasWFVars && !hasWFInit {
		issues = append(issues, Issue{
			Kind:    "mws",
			Level:   LevelWarning,
			Message: "wf.vars is used without wf.initVariables(); runtime environment may be uninitialized",
		})
	}

	if hasUtilsNamespace && !utilsArrayCallPattern.MatchString(code) {
		issues = append(issues, Issue{
			Kind:    "mws",
			Level:   LevelWarning,
			Message: "_utils namespace is referenced, but _utils.array() usage was not detected",
		})
	}

	if hasWFNamespace && !hasWFVars && !hasWFInit {
		issues = append(issues, Issue{
			Kind:    "mws",
			Level:   LevelWarning,
			Message: "wf namespace detected without wf.vars or wf.initVariables patterns; verify MWS API usage",
		})
	}

	return issues
}
