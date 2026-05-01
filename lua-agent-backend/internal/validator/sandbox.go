package validator

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// SandboxValidator executes Lua with a timeout and captured print output.
type SandboxValidator struct {
	timeout time.Duration
}

func NewSandboxValidator(timeout time.Duration) *SandboxValidator {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}

	return &SandboxValidator{timeout: timeout}
}

func (v *SandboxValidator) Validate(ctx context.Context, code string) *SandboxResult {
	runCtx, cancel := context.WithTimeout(ctx, v.timeout)
	defer cancel()

	var stdout bytes.Buffer
	L := lua.NewState(lua.Options{
		SkipOpenLibs: true,
	})
	defer L.Close()

	L.SetContext(runCtx)
	openSafeLibs(L)
	registerPrint(L, &stdout)
	registerMWSStubs(L)

	if _, err := L.LoadString(code); err != nil {
		return &SandboxResult{
			OK:     false,
			Stdout: stdout.String(),
			Error:  fmt.Sprintf("sandbox compile failed: %v", err),
		}
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- L.DoString(code)
	}()

	select {
	case <-runCtx.Done():
		return &SandboxResult{
			OK:       false,
			Stdout:   stdout.String(),
			Error:    runCtx.Err().Error(),
			TimedOut: true,
			Duration: v.timeout,
		}
	case err := <-errCh:
		return &SandboxResult{
			OK:       err == nil,
			Stdout:   stdout.String(),
			Duration: v.timeout,
			Error:    errorString(err),
		}
	}
}

func openSafeLibs(L *lua.LState) {
	lua.OpenBase(L)
	lua.OpenMath(L)
	lua.OpenString(L)
	lua.OpenTable(L)
}

func registerPrint(L *lua.LState, stdout *bytes.Buffer) {
	L.SetGlobal("print", L.NewFunction(func(state *lua.LState) int {
		top := state.GetTop()
		values := make([]string, 0, top)
		for i := 1; i <= top; i++ {
			values = append(values, state.CheckAny(i).String())
		}
		stdout.WriteString(strings.Join(values, "\t"))
		stdout.WriteByte('\n')
		return 0
	}))
}

func registerMWSStubs(L *lua.LState) {
	wf := L.NewTable()
	vars := L.NewTable()
	wf.RawSetString("vars", vars)
	wf.RawSetString("initVariables", L.NewFunction(func(state *lua.LState) int {
		state.Push(vars)
		return 1
	}))
	L.SetGlobal("wf", wf)

	utils := L.NewTable()
	utils.RawSetString("array", L.NewFunction(func(state *lua.LState) int {
		table := state.NewTable()
		top := state.GetTop()
		for i := 1; i <= top; i++ {
			table.Append(state.CheckAny(i))
		}
		state.Push(table)
		return 1
	}))
	L.SetGlobal("_utils", utils)
}

func errorString(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}
