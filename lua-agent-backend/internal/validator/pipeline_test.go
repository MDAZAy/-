package validator

import (
	"context"
	"testing"
	"time"
)

func TestValidatorRejectsForbiddenAPI(t *testing.T) {
	v := New(500 * time.Millisecond)

	result := v.Validate(context.Background(), `os.execute("rm -rf /")`)

	if result.OK {
		t.Fatal("expected validator to reject forbidden API")
	}
	if len(result.Issues) == 0 {
		t.Fatal("expected validation issues")
	}
}

func TestValidatorAcceptsSimpleLua(t *testing.T) {
	v := New(500 * time.Millisecond)

	result := v.Validate(context.Background(), `
local sum = 0
for _, value in ipairs({1, 2, 3}) do
  sum = sum + value
end
print(sum)
`)

	if !result.OK {
		t.Fatalf("expected valid Lua, got issues: %+v", result.Issues)
	}
	if result.Sandbox == nil || !result.Sandbox.OK {
		t.Fatalf("expected sandbox success, got: %+v", result.Sandbox)
	}
}

func TestValidatorWarnsOnWFVarsWithoutInit(t *testing.T) {
	v := New(500 * time.Millisecond)

	result := v.Validate(context.Background(), `wf.vars["name"] = "demo"`)

	if !result.OK {
		t.Fatalf("warning-only code should still be considered valid: %+v", result.Issues)
	}

	foundWarning := false
	for _, issue := range result.Issues {
		if issue.Kind == "mws" && issue.Level == LevelWarning {
			foundWarning = true
		}
	}

	if !foundWarning {
		t.Fatal("expected MWS warning")
	}
}
