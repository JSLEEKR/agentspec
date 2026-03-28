package spec

import (
	"strings"
	"testing"
)

func TestValidate_ValidSpec(t *testing.T) {
	s := &Spec{
		Name:  "valid",
		Input: Input{Message: "test"},
		Expect: Expect{
			Tools: []ToolExpect{{Name: "foo"}},
		},
	}
	if err := Validate(s); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_MissingName(t *testing.T) {
	s := &Spec{
		Input: Input{Message: "test"},
		Expect: Expect{
			Tools: []ToolExpect{{Name: "foo"}},
		},
	}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "name is required") {
		t.Errorf("error = %q, want name required", err.Error())
	}
}

func TestValidate_MissingMessage(t *testing.T) {
	s := &Spec{
		Name: "test",
		Expect: Expect{
			Tools: []ToolExpect{{Name: "foo"}},
		},
	}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "input.message is required") {
		t.Errorf("error = %q", err.Error())
	}
}

func TestValidate_EmptyExpect(t *testing.T) {
	s := &Spec{
		Name:   "test",
		Input:  Input{Message: "test"},
		Expect: Expect{},
	}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "expect must define at least one") {
		t.Errorf("error = %q", err.Error())
	}
}

func TestValidate_ToolMissingName(t *testing.T) {
	s := &Spec{
		Name:  "test",
		Input: Input{Message: "test"},
		Expect: Expect{
			Tools: []ToolExpect{{Name: ""}},
		},
	}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "tools[0].name is required") {
		t.Errorf("error = %q", err.Error())
	}
}

func TestValidate_InvalidRegexInArgs(t *testing.T) {
	s := &Spec{
		Name:  "test",
		Input: Input{Message: "test"},
		Expect: Expect{
			Tools: []ToolExpect{{
				Name: "foo",
				Args: map[string]interface{}{
					"key": map[string]interface{}{
						"regex": "[invalid",
					},
				},
			}},
		},
	}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestValidate_UnknownMatcherKey(t *testing.T) {
	s := &Spec{
		Name:  "test",
		Input: Input{Message: "test"},
		Expect: Expect{
			Tools: []ToolExpect{{
				Name: "foo",
				Args: map[string]interface{}{
					"key": map[string]interface{}{
						"badkey": "value",
					},
				},
			}},
		},
	}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected error for unknown matcher key")
	}
	if !strings.Contains(err.Error(), "unknown matcher key") {
		t.Errorf("error = %q", err.Error())
	}
}

func TestValidate_EmptyResponse(t *testing.T) {
	s := &Spec{
		Name:  "test",
		Input: Input{Message: "test"},
		Expect: Expect{
			Response: &Response{},
		},
	}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected error for empty response")
	}
}

func TestValidate_InvalidResponseRegex(t *testing.T) {
	s := &Spec{
		Name:  "test",
		Input: Input{Message: "test"},
		Expect: Expect{
			Response: &Response{Regex: "[bad"},
		},
	}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected error for invalid response regex")
	}
}

func TestValidate_ValidResponseContains(t *testing.T) {
	s := &Spec{
		Name:  "test",
		Input: Input{Message: "test"},
		Expect: Expect{
			Response: &Response{Contains: "hello"},
		},
	}
	if err := Validate(s); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_ValidResponseExact(t *testing.T) {
	s := &Spec{
		Name:  "test",
		Input: Input{Message: "test"},
		Expect: Expect{
			Response: &Response{Exact: "exact text"},
		},
	}
	if err := Validate(s); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_ConstraintNoField(t *testing.T) {
	s := &Spec{
		Name:  "test",
		Input: Input{Message: "test"},
		Expect: Expect{
			Constraints: []Constraint{{}},
		},
	}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected error for empty constraint")
	}
}

func TestValidate_ValidConstraints(t *testing.T) {
	s := &Spec{
		Name:  "test",
		Input: Input{Message: "test"},
		Expect: Expect{
			Constraints: []Constraint{
				{NoTool: "bad_tool"},
				{MaxTools: 5},
				{Ordered: true},
			},
		},
	}
	if err := Validate(s); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	s := &Spec{}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected error")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if len(ve.Errors) < 2 {
		t.Errorf("expected multiple errors, got %d", len(ve.Errors))
	}
}

func TestValidate_ValidMatcherContains(t *testing.T) {
	s := &Spec{
		Name:  "test",
		Input: Input{Message: "test"},
		Expect: Expect{
			Tools: []ToolExpect{{
				Name: "foo",
				Args: map[string]interface{}{
					"key": map[string]interface{}{
						"contains": "value",
					},
				},
			}},
		},
	}
	if err := Validate(s); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_ValidMatcherAny(t *testing.T) {
	s := &Spec{
		Name:  "test",
		Input: Input{Message: "test"},
		Expect: Expect{
			Tools: []ToolExpect{{
				Name: "foo",
				Args: map[string]interface{}{
					"key": map[string]interface{}{
						"any": true,
					},
				},
			}},
		},
	}
	if err := Validate(s); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_WhitespaceName(t *testing.T) {
	s := &Spec{
		Name:  "   ",
		Input: Input{Message: "test"},
		Expect: Expect{
			Tools: []ToolExpect{{Name: "foo"}},
		},
	}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected error for whitespace-only name")
	}
}
