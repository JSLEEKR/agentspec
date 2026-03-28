package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/JSLEEKR/agentspec/internal/runner"
)

// FormatTable writes results in human-readable table format.
func FormatTable(w io.Writer, rr *runner.RunResult) {
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "agentspec — Agent Behavior Testing")
	fmt.Fprintln(w, "===================================")
	fmt.Fprintln(w, "")

	totalChecks := 0
	passedChecks := 0

	for _, sr := range rr.Specs {
		label := sr.SpecPath
		if label == "" {
			label = sr.SpecName
		}
		fmt.Fprintln(w, label)

		for _, c := range sr.Checks {
			totalChecks++
			if c.Passed {
				passedChecks++
				fmt.Fprintf(w, "  %s %s\n", checkMark(), c.Message)
			} else {
				fmt.Fprintf(w, "  %s %s\n", crossMark(), c.Message)
			}
		}
		fmt.Fprintln(w, "")
	}

	failedChecks := totalChecks - passedChecks
	fmt.Fprintf(w, "Results: %d passed, %d failed (%d specs)\n",
		passedChecks, failedChecks, len(rr.Specs))
}

func checkMark() string {
	return "PASS"
}

func crossMark() string {
	return "FAIL"
}

// FormatTableString returns the table as a string.
func FormatTableString(rr *runner.RunResult) string {
	var sb strings.Builder
	FormatTable(&sb, rr)
	return sb.String()
}
