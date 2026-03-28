package matcher

import (
	"fmt"
	"regexp"
	"strings"
)

// MatchMode represents a matching strategy.
type MatchMode int

const (
	ModeExact MatchMode = iota
	ModeContains
	ModeRegex
	ModeSchema
	ModeAny
)

// MatchResult holds the outcome of a match operation.
type MatchResult struct {
	Matched bool
	Detail  string
}

// MatchValue checks if an actual value matches an expected value.
// Expected can be a string (exact match) or a map with matcher keys.
func MatchValue(expected, actual interface{}) MatchResult {
	// Handle matcher map: {contains: ..., regex: ..., type: ..., any: ...}
	if m, ok := expected.(map[string]interface{}); ok {
		return matchWithMap(m, actual)
	}

	// Exact match
	return matchExact(expected, actual)
}

func matchExact(expected, actual interface{}) MatchResult {
	es := fmt.Sprintf("%v", expected)
	as := fmt.Sprintf("%v", actual)
	if es == as {
		return MatchResult{Matched: true, Detail: fmt.Sprintf("exact match: %q", es)}
	}
	return MatchResult{Matched: false, Detail: fmt.Sprintf("expected %q, got %q", es, as)}
}

func matchWithMap(m map[string]interface{}, actual interface{}) MatchResult {
	// any: true matches anything
	if v, ok := m["any"]; ok {
		if b, ok := v.(bool); ok && b {
			return MatchResult{Matched: true, Detail: "any: matched"}
		}
	}

	as := fmt.Sprintf("%v", actual)

	// contains
	if v, ok := m["contains"]; ok {
		needle := fmt.Sprintf("%v", v)
		if strings.Contains(as, needle) {
			return MatchResult{Matched: true, Detail: fmt.Sprintf("contains %q", needle)}
		}
		return MatchResult{Matched: false, Detail: fmt.Sprintf("expected to contain %q, got %q", needle, as)}
	}

	// regex
	if v, ok := m["regex"]; ok {
		pattern := fmt.Sprintf("%v", v)
		re, err := regexp.Compile(pattern)
		if err != nil {
			return MatchResult{Matched: false, Detail: fmt.Sprintf("invalid regex %q: %v", pattern, err)}
		}
		if re.MatchString(as) {
			return MatchResult{Matched: true, Detail: fmt.Sprintf("regex /%s/ matched", pattern)}
		}
		return MatchResult{Matched: false, Detail: fmt.Sprintf("regex /%s/ did not match %q", pattern, as)}
	}

	// type (schema-style)
	if v, ok := m["type"]; ok {
		return matchType(fmt.Sprintf("%v", v), actual)
	}

	return MatchResult{Matched: false, Detail: "no known matcher in map"}
}

func matchType(expectedType string, actual interface{}) MatchResult {
	var actualType string
	switch actual.(type) {
	case string:
		actualType = "string"
	case float64, float32:
		actualType = "number"
	case int, int64, int32:
		actualType = "number"
	case bool:
		actualType = "boolean"
	case nil:
		actualType = "null"
	case []interface{}:
		actualType = "array"
	case map[string]interface{}:
		actualType = "object"
	default:
		actualType = fmt.Sprintf("%T", actual)
	}

	if actualType == expectedType {
		return MatchResult{Matched: true, Detail: fmt.Sprintf("type %s matched", expectedType)}
	}
	return MatchResult{Matched: false, Detail: fmt.Sprintf("expected type %s, got %s", expectedType, actualType)}
}
