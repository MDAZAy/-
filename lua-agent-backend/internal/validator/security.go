package validator

import (
	"fmt"
	"regexp"
)

type SecurityValidator struct {
	patterns []forbiddenPattern
}

type forbiddenPattern struct {
	name    string
	pattern *regexp.Regexp
}

func NewSecurityValidator() *SecurityValidator {
	return &SecurityValidator{
		patterns: []forbiddenPattern{
			{name: "os.execute", pattern: regexp.MustCompile(`\bos\.execute\s*\(`)},
			{name: "io.open", pattern: regexp.MustCompile(`\bio\.open\s*\(`)},
			{name: "dofile", pattern: regexp.MustCompile(`\bdofile\s*\(`)},
			{name: "loadfile", pattern: regexp.MustCompile(`\bloadfile\s*\(`)},
		},
	}
}

func (v *SecurityValidator) Validate(code string) []Issue {
	var issues []Issue
	for _, item := range v.patterns {
		if item.pattern.MatchString(code) {
			issues = append(issues, Issue{
				Kind:    "security",
				Level:   LevelError,
				Message: fmt.Sprintf("forbidden lua API detected: %s", item.name),
			})
		}
	}

	return issues
}
