package reporter

import "github.com/JSLEEKR/agentspec/internal/runner"

// Summary holds aggregated test results.
type Summary struct {
	Total  int `json:"total"`
	Passed int `json:"passed"`
	Failed int `json:"failed"`
	Checks int `json:"checks"`
}

// BuildSummary creates a summary from run results.
// Total, Passed, and Failed all count specs. Checks counts individual checks.
func BuildSummary(rr *runner.RunResult) Summary {
	checks := 0
	for _, sr := range rr.Specs {
		checks += len(sr.Checks)
	}
	return Summary{
		Total:  len(rr.Specs),
		Passed: rr.TotalPassed,
		Failed: rr.TotalFailed,
		Checks: checks,
	}
}
