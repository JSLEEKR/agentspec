package reporter

import (
	"strings"
	"testing"

	"github.com/JSLEEKR/agentspec/internal/matcher"
	"github.com/JSLEEKR/agentspec/internal/runner"
)

func TestFormatTable_AllPass(t *testing.T) {
	rr := &runner.RunResult{
		Specs: []runner.SpecResult{
			{
				SpecName: "test spec",
				SpecPath: "specs/test.yaml",
				Passed:   true,
				Checks: []matcher.CheckResult{
					{Passed: true, Message: "Tool call: read_file"},
					{Passed: true, Message: "Response contains \"hello\""},
				},
			},
		},
		TotalPassed: 1,
		TotalFailed: 0,
	}

	output := FormatTableString(rr)

	if !strings.Contains(output, "agentspec") {
		t.Error("output should contain header")
	}
	if !strings.Contains(output, "specs/test.yaml") {
		t.Error("output should contain spec path")
	}
	if !strings.Contains(output, "PASS") {
		t.Error("output should contain PASS marker")
	}
	if !strings.Contains(output, "2 passed, 0 failed") {
		t.Errorf("output should contain summary, got:\n%s", output)
	}
}

func TestFormatTable_WithFailures(t *testing.T) {
	rr := &runner.RunResult{
		Specs: []runner.SpecResult{
			{
				SpecName: "fail spec",
				Passed:   false,
				Checks: []matcher.CheckResult{
					{Passed: true, Message: "Tool call: search"},
					{Passed: false, Message: "Tool call: summarize -- not called"},
				},
			},
		},
		TotalPassed: 0,
		TotalFailed: 1,
	}

	output := FormatTableString(rr)

	if !strings.Contains(output, "FAIL") {
		t.Error("output should contain FAIL marker")
	}
	if !strings.Contains(output, "1 passed, 1 failed") {
		t.Errorf("unexpected summary: %s", output)
	}
}

func TestFormatTable_MultipleSpecs(t *testing.T) {
	rr := &runner.RunResult{
		Specs: []runner.SpecResult{
			{SpecName: "spec1", Passed: true, Checks: []matcher.CheckResult{{Passed: true, Message: "ok"}}},
			{SpecName: "spec2", Passed: false, Checks: []matcher.CheckResult{{Passed: false, Message: "fail"}}},
		},
		TotalPassed: 1,
		TotalFailed: 1,
	}

	output := FormatTableString(rr)
	if !strings.Contains(output, "2 specs") {
		t.Errorf("should show 2 specs: %s", output)
	}
}

func TestFormatTable_UsesPathOverName(t *testing.T) {
	rr := &runner.RunResult{
		Specs: []runner.SpecResult{
			{SpecName: "test", SpecPath: "my/path.yaml", Passed: true, Checks: []matcher.CheckResult{}},
		},
	}
	output := FormatTableString(rr)
	if !strings.Contains(output, "my/path.yaml") {
		t.Errorf("should use path: %s", output)
	}
}

func TestFormatTable_FallsBackToName(t *testing.T) {
	rr := &runner.RunResult{
		Specs: []runner.SpecResult{
			{SpecName: "my test", Passed: true, Checks: []matcher.CheckResult{}},
		},
	}
	output := FormatTableString(rr)
	if !strings.Contains(output, "my test") {
		t.Errorf("should use name: %s", output)
	}
}

func TestFormatTable_NoChecks(t *testing.T) {
	rr := &runner.RunResult{
		Specs: []runner.SpecResult{
			{SpecName: "empty", Passed: true, Checks: []matcher.CheckResult{}},
		},
	}
	output := FormatTableString(rr)
	if !strings.Contains(output, "0 passed, 0 failed") {
		t.Errorf("unexpected: %s", output)
	}
}
