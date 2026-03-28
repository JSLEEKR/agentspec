package reporter

import "github.com/JSLEEKR/agentspec/internal/runner"

// Summary holds aggregated test results.
type Summary struct {
	Total  int `json:"total"`
	Passed int `json:"passed"`
	Failed int `json:"failed"`
	Specs  int `json:"specs"`
}

// BuildSummary creates a summary from run results.
func BuildSummary(rr *runner.RunResult) Summary {
	total := 0
	for _, sr := range rr.Specs {
		total += len(sr.Checks)
	}
	return Summary{
		Total:  total,
		Passed: rr.TotalPassed,
		Failed: rr.TotalFailed,
		Specs:  len(rr.Specs),
	}
}
