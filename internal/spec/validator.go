package spec

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError holds all validation issues for a spec.
type ValidationError struct {
	Errors []string
}

func (ve *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s", strings.Join(ve.Errors, "; "))
}

// Validate checks a spec for structural correctness.
func Validate(s *Spec) error {
	var errs []string

	if strings.TrimSpace(s.Name) == "" {
		errs = append(errs, "spec name is required")
	}
	if strings.TrimSpace(s.Input.Message) == "" {
		errs = append(errs, "input.message is required")
	}
	if len(s.Expect.Tools) == 0 && s.Expect.Response == nil && len(s.Expect.Constraints) == 0 {
		errs = append(errs, "expect must define at least one of: tools, response, constraints")
	}

	for i, t := range s.Expect.Tools {
		if strings.TrimSpace(t.Name) == "" {
			errs = append(errs, fmt.Sprintf("expect.tools[%d].name is required", i))
		}
		// Validate arg matchers
		for key, val := range t.Args {
			if m, ok := val.(map[string]interface{}); ok {
				if err := validateArgMatcher(i, key, m); err != "" {
					errs = append(errs, err)
				}
			}
		}
	}

	if s.Expect.Response != nil {
		if err := validateResponse(s.Expect.Response); err != "" {
			errs = append(errs, err)
		}
	}

	for i, c := range s.Expect.Constraints {
		if err := validateConstraint(i, c); err != "" {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return &ValidationError{Errors: errs}
	}
	return nil
}

func validateArgMatcher(toolIdx int, key string, m map[string]interface{}) string {
	validKeys := map[string]bool{
		"contains": true, "regex": true, "type": true, "any": true,
	}
	for k := range m {
		if !validKeys[k] {
			return fmt.Sprintf("expect.tools[%d].args.%s has unknown matcher key %q", toolIdx, key, k)
		}
	}
	if r, ok := m["regex"]; ok {
		if rs, ok := r.(string); ok {
			if _, err := regexp.Compile(rs); err != nil {
				return fmt.Sprintf("expect.tools[%d].args.%s.regex is invalid: %v", toolIdx, key, err)
			}
		}
	}
	return ""
}

func validateResponse(r *Response) string {
	if r.Contains == "" && r.Exact == "" && r.Regex == "" {
		return "expect.response must have at least one of: contains, exact, regex"
	}
	if r.Regex != "" {
		if _, err := regexp.Compile(r.Regex); err != nil {
			return fmt.Sprintf("expect.response.regex is invalid: %v", err)
		}
	}
	return ""
}

func validateConstraint(idx int, c Constraint) string {
	if c.MaxTools < 0 {
		return fmt.Sprintf("expect.constraints[%d].max_tools must be positive", idx)
	}
	hasField := c.NoTool != "" || c.MaxTools != 0 || c.Ordered
	if !hasField {
		return fmt.Sprintf("expect.constraints[%d] has no valid constraint field", idx)
	}
	return ""
}
