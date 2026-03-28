package reporter

import (
	"testing"

	"github.com/JSLEEKR/agentspec/internal/matcher"
	"github.com/JSLEEKR/agentspec/internal/runner"
)

func TestBuildSummary_Basic(t *testing.T) {
	rr := &runner.RunResult{
		Specs: []runner.SpecResult{
			{Passed: true, Checks: []matcher.CheckResult{{Passed: true}, {Passed: true}}},
			{Passed: false, Checks: []matcher.CheckResult{{Passed: false}}},
		},
		TotalPassed: 1,
		TotalFailed: 1,
	}

	s := BuildSummary(rr)
	if s.Total != 3 {
		t.Errorf("total = %d, want 3", s.Total)
	}
	if s.Passed != 1 {
		t.Errorf("passed = %d, want 1", s.Passed)
	}
	if s.Failed != 1 {
		t.Errorf("failed = %d, want 1", s.Failed)
	}
	if s.Specs != 2 {
		t.Errorf("specs = %d, want 2", s.Specs)
	}
}

func TestBuildSummary_Empty(t *testing.T) {
	rr := &runner.RunResult{}
	s := BuildSummary(rr)
	if s.Total != 0 || s.Passed != 0 || s.Failed != 0 || s.Specs != 0 {
		t.Errorf("expected all zeros: %+v", s)
	}
}

func TestBuildSummary_AllPass(t *testing.T) {
	rr := &runner.RunResult{
		Specs: []runner.SpecResult{
			{Passed: true, Checks: []matcher.CheckResult{{Passed: true}}},
			{Passed: true, Checks: []matcher.CheckResult{{Passed: true}}},
		},
		TotalPassed: 2,
	}
	s := BuildSummary(rr)
	if s.Failed != 0 {
		t.Errorf("failed = %d, want 0", s.Failed)
	}
}
