package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/JSLEEKR/agentspec/internal/runner"
)

// JSONReport represents the structured JSON output.
type JSONReport struct {
	Summary Summary      `json:"summary"`
	Specs   []JSONSpec   `json:"specs"`
}

// JSONSpec represents a single spec result in JSON.
type JSONSpec struct {
	Name    string      `json:"name"`
	Path    string      `json:"path,omitempty"`
	Passed  bool        `json:"passed"`
	Checks  []JSONCheck `json:"checks"`
}

// JSONCheck represents a single check result in JSON.
type JSONCheck struct {
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
}

// FormatJSON writes results in JSON format.
func FormatJSON(w io.Writer, rr *runner.RunResult) error {
	report := JSONReport{
		Summary: BuildSummary(rr),
	}

	for _, sr := range rr.Specs {
		js := JSONSpec{
			Name:   sr.SpecName,
			Path:   sr.SpecPath,
			Passed: sr.Passed,
		}
		for _, c := range sr.Checks {
			js.Checks = append(js.Checks, JSONCheck{
				Passed:  c.Passed,
				Message: c.Message,
			})
		}
		report.Specs = append(report.Specs, js)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		return fmt.Errorf("encode JSON report: %w", err)
	}
	return nil
}
