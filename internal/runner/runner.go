package runner

import (
	"fmt"

	"github.com/JSLEEKR/agentspec/internal/loader"
	"github.com/JSLEEKR/agentspec/internal/matcher"
	"github.com/JSLEEKR/agentspec/internal/spec"
)

// SpecResult holds the results of running a single spec.
type SpecResult struct {
	SpecName string
	SpecPath string
	Checks   []matcher.CheckResult
	Passed   bool
}

// RunResult holds the results of running all specs.
type RunResult struct {
	Specs       []SpecResult
	TotalPassed int
	TotalFailed int
}

// Run executes a single spec against an execution log.
func Run(s *spec.Spec, log *loader.ExecutionLog) SpecResult {
	var checks []matcher.CheckResult

	// Match tools
	if len(s.Expect.Tools) > 0 {
		checks = append(checks, matcher.MatchTools(s.Expect.Tools, log.ToolCalls)...)
	}

	// Match response
	if s.Expect.Response != nil {
		checks = append(checks, matcher.MatchResponse(s.Expect.Response, log.Response)...)
	}

	// Match constraints
	if len(s.Expect.Constraints) > 0 {
		checks = append(checks, matcher.MatchConstraints(s.Expect.Constraints, s.Expect.Tools, log.ToolCalls)...)
	}

	allPassed := true
	for _, c := range checks {
		if !c.Passed {
			allPassed = false
			break
		}
	}

	return SpecResult{
		SpecName: s.Name,
		Checks:   checks,
		Passed:   allPassed,
	}
}

// RunAll executes multiple specs against their corresponding logs.
func RunAll(specs []*spec.Spec, paths []string, logs []*loader.ExecutionLog) (*RunResult, error) {
	if len(specs) != len(logs) {
		return nil, fmt.Errorf("spec count (%d) does not match log count (%d)", len(specs), len(logs))
	}

	result := &RunResult{}
	for i, s := range specs {
		sr := Run(s, logs[i])
		if len(paths) > i {
			sr.SpecPath = paths[i]
		}
		result.Specs = append(result.Specs, sr)
		if sr.Passed {
			result.TotalPassed++
		} else {
			result.TotalFailed++
		}
	}

	return result, nil
}

// RunSingle runs a spec against a single log (convenience wrapper).
func RunSingle(s *spec.Spec, log *loader.ExecutionLog) *RunResult {
	sr := Run(s, log)
	result := &RunResult{
		Specs: []SpecResult{sr},
	}
	if sr.Passed {
		result.TotalPassed = 1
	} else {
		result.TotalFailed = 1
	}
	return result
}
