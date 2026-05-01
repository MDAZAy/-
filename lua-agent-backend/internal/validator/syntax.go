package validator

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// SyntaxValidator compiles Lua code to ensure it is parsable by gopher-lua.
type SyntaxValidator struct{}

func NewSyntaxValidator() *SyntaxValidator {
	return &SyntaxValidator{}
}

func (v *SyntaxValidator) Validate(code string) []Issue {
	L := lua.NewState(lua.Options{
		SkipOpenLibs: true,
	})
	defer L.Close()

	if _, err := L.LoadString(code); err != nil {
		return []Issue{
			{
				Kind:    "syntax",
				Level:   LevelError,
				Message: fmt.Sprintf("lua syntax check failed: %v", err),
			},
		}
	}

	return nil
}
