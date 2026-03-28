package reporter

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/JSLEEKR/agentspec/internal/matcher"
	"github.com/JSLEEKR/agentspec/internal/runner"
)

func TestFormatJSON_Valid(t *testing.T) {
	rr := &runner.RunResult{
		Specs: []runner.SpecResult{
			{
				SpecName: "test",
				SpecPath: "specs/test.yaml",
				Passed:   true,
				Checks: []matcher.CheckResult{
					{Passed: true, Message: "Tool call: read_file"},
				},
			},
		},
		TotalPassed: 1,
		TotalFailed: 0,
	}

	var sb strings.Builder
	if err := FormatJSON(&sb, rr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var report JSONReport
	if err := json.Unmarshal([]byte(sb.String()), &report); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if report.Summary.Passed != 1 {
		t.Errorf("summary.passed = %d, want 1", report.Summary.Passed)
	}
	if len(report.Specs) != 1 {
		t.Fatalf("specs = %d, want 1", len(report.Specs))
	}
	if report.Specs[0].Name != "test" {
		t.Errorf("spec name = %q", report.Specs[0].Name)
	}
	if !report.Specs[0].Passed {
		t.Error("spec should pass")
	}
}

func TestFormatJSON_WithFailures(t *testing.T) {
	rr := &runner.RunResult{
		Specs: []runner.SpecResult{
			{
				SpecName: "fail",
				Passed:   false,
				Checks: []matcher.CheckResult{
					{Passed: false, Message: "not found"},
				},
			},
		},
		TotalPassed: 0,
		TotalFailed: 1,
	}

	var sb strings.Builder
	if err := FormatJSON(&sb, rr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var report JSONReport
	if err := json.Unmarshal([]byte(sb.String()), &report); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if report.Summary.Failed != 1 {
		t.Errorf("summary.failed = %d, want 1", report.Summary.Failed)
	}
	if report.Specs[0].Passed {
		t.Error("spec should fail")
	}
}

func TestFormatJSON_MultipleSpecs(t *testing.T) {
	rr := &runner.RunResult{
		Specs: []runner.SpecResult{
			{SpecName: "a", Passed: true, Checks: []matcher.CheckResult{{Passed: true, Message: "ok"}}},
			{SpecName: "b", Passed: false, Checks: []matcher.CheckResult{{Passed: false, Message: "fail"}}},
		},
		TotalPassed: 1,
		TotalFailed: 1,
	}

	var sb strings.Builder
	if err := FormatJSON(&sb, rr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var report JSONReport
	if err := json.Unmarshal([]byte(sb.String()), &report); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if report.Summary.Specs != 2 {
		t.Errorf("specs = %d, want 2", report.Summary.Specs)
	}
}

func TestFormatJSON_EmptyResults(t *testing.T) {
	rr := &runner.RunResult{}

	var sb strings.Builder
	if err := FormatJSON(&sb, rr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var report JSONReport
	if err := json.Unmarshal([]byte(sb.String()), &report); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if report.Summary.Total != 0 {
		t.Errorf("total = %d, want 0", report.Summary.Total)
	}
}

func TestFormatJSON_CheckDetails(t *testing.T) {
	rr := &runner.RunResult{
		Specs: []runner.SpecResult{
			{
				SpecName: "detail",
				Passed:   true,
				Checks: []matcher.CheckResult{
					{Passed: true, Message: "check 1"},
					{Passed: true, Message: "check 2"},
				},
			},
		},
		TotalPassed: 1,
	}

	var sb strings.Builder
	FormatJSON(&sb, rr)

	var report JSONReport
	json.Unmarshal([]byte(sb.String()), &report)

	if len(report.Specs[0].Checks) != 2 {
		t.Errorf("checks = %d, want 2", len(report.Specs[0].Checks))
	}
}
