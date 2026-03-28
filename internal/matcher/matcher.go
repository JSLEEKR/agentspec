package matcher

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/JSLEEKR/agentspec/internal/loader"
	"github.com/JSLEEKR/agentspec/internal/spec"
)

// CheckResult represents the result of checking a single expectation.
type CheckResult struct {
	Passed  bool
	Message string
}

// MatchTools checks that expected tool calls are present in the execution log.
func MatchTools(expected []spec.ToolExpect, actual []loader.ToolCall) []CheckResult {
	var results []CheckResult

	for _, exp := range expected {
		found := false
		for _, act := range actual {
			if act.Name != exp.Name {
				continue
			}
			// Tool name matches, check args
			if len(exp.Args) == 0 {
				found = true
				results = append(results, CheckResult{
					Passed:  true,
					Message: fmt.Sprintf("Tool call: %s", exp.Name),
				})
				break
			}

			allArgsMatch := true
			var argDetails []string
			for key, expVal := range exp.Args {
				actVal, exists := act.Arguments[key]
				if !exists {
					allArgsMatch = false
					argDetails = append(argDetails, fmt.Sprintf("%s: missing", key))
					continue
				}
				mr := MatchValue(expVal, actVal)
				if !mr.Matched {
					allArgsMatch = false
					argDetails = append(argDetails, fmt.Sprintf("%s: %s", key, mr.Detail))
				} else {
					argDetails = append(argDetails, fmt.Sprintf("%s=%v", key, actVal))
				}
			}

			if allArgsMatch {
				found = true
				results = append(results, CheckResult{
					Passed:  true,
					Message: fmt.Sprintf("Tool call: %s (%s)", exp.Name, strings.Join(argDetails, ", ")),
				})
				break
			}
		}

		if !found {
			detail := exp.Name
			if len(exp.Args) > 0 {
				var parts []string
				for k, v := range exp.Args {
					parts = append(parts, fmt.Sprintf("%s=%v", k, v))
				}
				detail += " (" + strings.Join(parts, ", ") + ")"
			}
			results = append(results, CheckResult{
				Passed:  false,
				Message: fmt.Sprintf("Tool call: %s -- expected but not called", detail),
			})
		}
	}

	return results
}

// MatchResponse checks if the agent's response matches expectations.
func MatchResponse(expected *spec.Response, actual string) []CheckResult {
	if expected == nil {
		return nil
	}

	var results []CheckResult

	if expected.Contains != "" {
		if strings.Contains(actual, expected.Contains) {
			results = append(results, CheckResult{
				Passed:  true,
				Message: fmt.Sprintf("Response contains %q", expected.Contains),
			})
		} else {
			results = append(results, CheckResult{
				Passed:  false,
				Message: fmt.Sprintf("Response does not contain %q", expected.Contains),
			})
		}
	}

	if expected.Exact != "" {
		if actual == expected.Exact {
			results = append(results, CheckResult{
				Passed:  true,
				Message: "Response exact match",
			})
		} else {
			results = append(results, CheckResult{
				Passed:  false,
				Message: fmt.Sprintf("Response exact mismatch: expected %q", expected.Exact),
			})
		}
	}

	if expected.Regex != "" {
		re, err := regexp.Compile(expected.Regex)
		if err != nil {
			results = append(results, CheckResult{
				Passed:  false,
				Message: fmt.Sprintf("Response regex invalid: %v", err),
			})
		} else if re.MatchString(actual) {
			results = append(results, CheckResult{
				Passed:  true,
				Message: fmt.Sprintf("Response matches /%s/", expected.Regex),
			})
		} else {
			results = append(results, CheckResult{
				Passed:  false,
				Message: fmt.Sprintf("Response does not match /%s/", expected.Regex),
			})
		}
	}

	return results
}

// MatchConstraints checks behavioral constraints against actual tool calls.
func MatchConstraints(constraints []spec.Constraint, expectedTools []spec.ToolExpect, actual []loader.ToolCall) []CheckResult {
	var results []CheckResult

	for _, c := range constraints {
		if c.NoTool != "" {
			found := false
			for _, tc := range actual {
				if tc.Name == c.NoTool {
					found = true
					break
				}
			}
			if found {
				results = append(results, CheckResult{
					Passed:  false,
					Message: fmt.Sprintf("Constraint: no_%s calls -- violated, %s was called", c.NoTool, c.NoTool),
				})
			} else {
				results = append(results, CheckResult{
					Passed:  true,
					Message: fmt.Sprintf("Constraint: no %s calls", c.NoTool),
				})
			}
		}

		if c.MaxTools > 0 {
			if len(actual) <= c.MaxTools {
				results = append(results, CheckResult{
					Passed:  true,
					Message: fmt.Sprintf("Constraint: max %d tool calls (%d used)", c.MaxTools, len(actual)),
				})
			} else {
				results = append(results, CheckResult{
					Passed:  false,
					Message: fmt.Sprintf("Constraint: max %d tool calls -- exceeded (%d used)", c.MaxTools, len(actual)),
				})
			}
		}

		if c.Ordered {
			results = append(results, checkOrdered(expectedTools, actual))
		}
	}

	return results
}

func checkOrdered(expected []spec.ToolExpect, actual []loader.ToolCall) CheckResult {
	if len(expected) == 0 {
		return CheckResult{Passed: true, Message: "Constraint: ordered execution (no tools expected)"}
	}

	lastIdx := -1
	for _, exp := range expected {
		found := false
		for i, act := range actual {
			if act.Name == exp.Name && i > lastIdx {
				lastIdx = i
				found = true
				break
			}
		}
		if !found {
			// Tool not found at all -- ordering check is N/A but we still report
			return CheckResult{
				Passed:  false,
				Message: fmt.Sprintf("Constraint: ordered execution -- %s not found after index %d", exp.Name, lastIdx),
			}
		}
	}

	return CheckResult{Passed: true, Message: "Constraint: ordered execution"}
}
