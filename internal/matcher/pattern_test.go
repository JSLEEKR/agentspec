package matcher

import (
	"testing"
)

func TestMatchValue_ExactString(t *testing.T) {
	r := MatchValue("main.go", "main.go")
	if !r.Matched {
		t.Errorf("expected match: %s", r.Detail)
	}
}

func TestMatchValue_ExactMismatch(t *testing.T) {
	r := MatchValue("main.go", "test.go")
	if r.Matched {
		t.Error("expected no match")
	}
}

func TestMatchValue_ExactNumber(t *testing.T) {
	r := MatchValue(42, 42)
	if !r.Matched {
		t.Errorf("expected match: %s", r.Detail)
	}
}

func TestMatchValue_Contains(t *testing.T) {
	m := map[string]interface{}{"contains": "weather"}
	r := MatchValue(m, "what is the weather in tokyo")
	if !r.Matched {
		t.Errorf("expected match: %s", r.Detail)
	}
}

func TestMatchValue_ContainsMismatch(t *testing.T) {
	m := map[string]interface{}{"contains": "weather"}
	r := MatchValue(m, "hello world")
	if r.Matched {
		t.Error("expected no match")
	}
}

func TestMatchValue_Regex(t *testing.T) {
	m := map[string]interface{}{"regex": `.*\.go$`}
	r := MatchValue(m, "main.go")
	if !r.Matched {
		t.Errorf("expected match: %s", r.Detail)
	}
}

func TestMatchValue_RegexMismatch(t *testing.T) {
	m := map[string]interface{}{"regex": `.*\.py$`}
	r := MatchValue(m, "main.go")
	if r.Matched {
		t.Error("expected no match")
	}
}

func TestMatchValue_RegexInvalid(t *testing.T) {
	m := map[string]interface{}{"regex": "[invalid"}
	r := MatchValue(m, "test")
	if r.Matched {
		t.Error("expected no match for invalid regex")
	}
}

func TestMatchValue_Any(t *testing.T) {
	m := map[string]interface{}{"any": true}
	r := MatchValue(m, "anything goes")
	if !r.Matched {
		t.Errorf("expected match: %s", r.Detail)
	}
}

func TestMatchValue_AnyFalse(t *testing.T) {
	m := map[string]interface{}{"any": false}
	r := MatchValue(m, "test")
	// any: false should not match via the any path
	if r.Matched {
		t.Error("any: false should not auto-match")
	}
}

func TestMatchValue_TypeString(t *testing.T) {
	m := map[string]interface{}{"type": "string"}
	r := MatchValue(m, "hello")
	if !r.Matched {
		t.Errorf("expected match: %s", r.Detail)
	}
}

func TestMatchValue_TypeNumber(t *testing.T) {
	m := map[string]interface{}{"type": "number"}
	r := MatchValue(m, float64(42))
	if !r.Matched {
		t.Errorf("expected match: %s", r.Detail)
	}
}

func TestMatchValue_TypeBoolean(t *testing.T) {
	m := map[string]interface{}{"type": "boolean"}
	r := MatchValue(m, true)
	if !r.Matched {
		t.Errorf("expected match: %s", r.Detail)
	}
}

func TestMatchValue_TypeMismatch(t *testing.T) {
	m := map[string]interface{}{"type": "number"}
	r := MatchValue(m, "not a number")
	if r.Matched {
		t.Error("expected no match")
	}
}

func TestMatchValue_TypeNull(t *testing.T) {
	m := map[string]interface{}{"type": "null"}
	r := MatchValue(m, nil)
	if !r.Matched {
		t.Errorf("expected match: %s", r.Detail)
	}
}

func TestMatchValue_TypeArray(t *testing.T) {
	m := map[string]interface{}{"type": "array"}
	r := MatchValue(m, []interface{}{1, 2, 3})
	if !r.Matched {
		t.Errorf("expected match: %s", r.Detail)
	}
}

func TestMatchValue_TypeObject(t *testing.T) {
	m := map[string]interface{}{"type": "object"}
	r := MatchValue(m, map[string]interface{}{"key": "val"})
	if !r.Matched {
		t.Errorf("expected match: %s", r.Detail)
	}
}

func TestMatchValue_NoMatcher(t *testing.T) {
	m := map[string]interface{}{}
	r := MatchValue(m, "test")
	if r.Matched {
		t.Error("expected no match for empty matcher map")
	}
}

func TestMatchValue_ExactEmptyString(t *testing.T) {
	r := MatchValue("", "")
	if !r.Matched {
		t.Errorf("expected match for empty strings: %s", r.Detail)
	}
}

func TestMatchValue_ContainsEmpty(t *testing.T) {
	m := map[string]interface{}{"contains": ""}
	r := MatchValue(m, "anything")
	if !r.Matched {
		t.Errorf("empty contains should match everything: %s", r.Detail)
	}
}

func TestMatchValue_RegexFullMatch(t *testing.T) {
	m := map[string]interface{}{"regex": `^exact$`}
	r := MatchValue(m, "exact")
	if !r.Matched {
		t.Errorf("expected match: %s", r.Detail)
	}
}

func TestMatchValue_RegexPartialMatch(t *testing.T) {
	m := map[string]interface{}{"regex": `\d+`}
	r := MatchValue(m, "abc123def")
	if !r.Matched {
		t.Errorf("expected match: %s", r.Detail)
	}
}
