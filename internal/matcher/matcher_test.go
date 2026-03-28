package matcher

import (
	"testing"

	"github.com/JSLEEKR/agentspec/internal/loader"
	"github.com/JSLEEKR/agentspec/internal/spec"
)

func TestMatchTools_SingleExact(t *testing.T) {
	expected := []spec.ToolExpect{{
		Name: "read_file",
		Args: map[string]interface{}{"path": "main.go"},
	}}
	actual := []loader.ToolCall{{
		Name:      "read_file",
		Arguments: map[string]interface{}{"path": "main.go"},
	}}
	results := MatchTools(expected, actual)
	if len(results) != 1 {
		t.Fatalf("results = %d, want 1", len(results))
	}
	if !results[0].Passed {
		t.Errorf("expected pass: %s", results[0].Message)
	}
}

func TestMatchTools_NotFound(t *testing.T) {
	expected := []spec.ToolExpect{{Name: "missing_tool"}}
	actual := []loader.ToolCall{{Name: "other_tool", Arguments: map[string]interface{}{}}}
	results := MatchTools(expected, actual)
	if len(results) != 1 {
		t.Fatalf("results = %d, want 1", len(results))
	}
	if results[0].Passed {
		t.Error("expected fail")
	}
}

func TestMatchTools_NoArgs(t *testing.T) {
	expected := []spec.ToolExpect{{Name: "list_files"}}
	actual := []loader.ToolCall{{Name: "list_files", Arguments: map[string]interface{}{}}}
	results := MatchTools(expected, actual)
	if !results[0].Passed {
		t.Errorf("expected pass: %s", results[0].Message)
	}
}

func TestMatchTools_ArgMismatch(t *testing.T) {
	expected := []spec.ToolExpect{{
		Name: "read_file",
		Args: map[string]interface{}{"path": "main.go"},
	}}
	actual := []loader.ToolCall{{
		Name:      "read_file",
		Arguments: map[string]interface{}{"path": "test.go"},
	}}
	results := MatchTools(expected, actual)
	if results[0].Passed {
		t.Error("expected fail for arg mismatch")
	}
}

func TestMatchTools_MissingArg(t *testing.T) {
	expected := []spec.ToolExpect{{
		Name: "tool",
		Args: map[string]interface{}{"key": "value"},
	}}
	actual := []loader.ToolCall{{
		Name:      "tool",
		Arguments: map[string]interface{}{},
	}}
	results := MatchTools(expected, actual)
	if results[0].Passed {
		t.Error("expected fail for missing arg")
	}
}

func TestMatchTools_ContainsArg(t *testing.T) {
	expected := []spec.ToolExpect{{
		Name: "search",
		Args: map[string]interface{}{
			"query": map[string]interface{}{"contains": "weather"},
		},
	}}
	actual := []loader.ToolCall{{
		Name:      "search",
		Arguments: map[string]interface{}{"query": "what is the weather today"},
	}}
	results := MatchTools(expected, actual)
	if !results[0].Passed {
		t.Errorf("expected pass: %s", results[0].Message)
	}
}

func TestMatchTools_RegexArg(t *testing.T) {
	expected := []spec.ToolExpect{{
		Name: "search",
		Args: map[string]interface{}{
			"query": map[string]interface{}{"regex": `weather.*tokyo`},
		},
	}}
	actual := []loader.ToolCall{{
		Name:      "search",
		Arguments: map[string]interface{}{"query": "weather in tokyo"},
	}}
	results := MatchTools(expected, actual)
	if !results[0].Passed {
		t.Errorf("expected pass: %s", results[0].Message)
	}
}

func TestMatchTools_MultipleExpected(t *testing.T) {
	expected := []spec.ToolExpect{
		{Name: "search"},
		{Name: "summarize"},
	}
	actual := []loader.ToolCall{
		{Name: "search", Arguments: map[string]interface{}{}},
		{Name: "summarize", Arguments: map[string]interface{}{}},
	}
	results := MatchTools(expected, actual)
	if len(results) != 2 {
		t.Fatalf("results = %d, want 2", len(results))
	}
	for i, r := range results {
		if !r.Passed {
			t.Errorf("result[%d] expected pass: %s", i, r.Message)
		}
	}
}

func TestMatchTools_PartialMatch(t *testing.T) {
	expected := []spec.ToolExpect{
		{Name: "search"},
		{Name: "missing"},
	}
	actual := []loader.ToolCall{
		{Name: "search", Arguments: map[string]interface{}{}},
	}
	results := MatchTools(expected, actual)
	if results[0].Passed != true {
		t.Error("search should pass")
	}
	if results[1].Passed != false {
		t.Error("missing should fail")
	}
}

func TestMatchResponse_Contains(t *testing.T) {
	resp := &spec.Response{Contains: "package main"}
	results := MatchResponse(resp, "Here is the file: package main...")
	if len(results) != 1 {
		t.Fatalf("results = %d, want 1", len(results))
	}
	if !results[0].Passed {
		t.Errorf("expected pass: %s", results[0].Message)
	}
}

func TestMatchResponse_ContainsMismatch(t *testing.T) {
	resp := &spec.Response{Contains: "not here"}
	results := MatchResponse(resp, "something else")
	if results[0].Passed {
		t.Error("expected fail")
	}
}

func TestMatchResponse_Exact(t *testing.T) {
	resp := &spec.Response{Exact: "Done."}
	results := MatchResponse(resp, "Done.")
	if !results[0].Passed {
		t.Errorf("expected pass: %s", results[0].Message)
	}
}

func TestMatchResponse_ExactMismatch(t *testing.T) {
	resp := &spec.Response{Exact: "Done."}
	results := MatchResponse(resp, "Done!")
	if results[0].Passed {
		t.Error("expected fail")
	}
}

func TestMatchResponse_Regex(t *testing.T) {
	resp := &spec.Response{Regex: `\d+ results`}
	results := MatchResponse(resp, "Found 42 results")
	if !results[0].Passed {
		t.Errorf("expected pass: %s", results[0].Message)
	}
}

func TestMatchResponse_RegexMismatch(t *testing.T) {
	resp := &spec.Response{Regex: `^\d+$`}
	results := MatchResponse(resp, "not a number")
	if results[0].Passed {
		t.Error("expected fail")
	}
}

func TestMatchResponse_InvalidRegex(t *testing.T) {
	resp := &spec.Response{Regex: "[bad"}
	results := MatchResponse(resp, "test")
	if results[0].Passed {
		t.Error("expected fail for invalid regex")
	}
}

func TestMatchResponse_Nil(t *testing.T) {
	results := MatchResponse(nil, "anything")
	if len(results) != 0 {
		t.Errorf("results = %d, want 0", len(results))
	}
}

func TestMatchConstraints_NoTool_Pass(t *testing.T) {
	constraints := []spec.Constraint{{NoTool: "write_file"}}
	actual := []loader.ToolCall{{Name: "read_file", Arguments: map[string]interface{}{}}}
	results := MatchConstraints(constraints, nil, actual)
	if !results[0].Passed {
		t.Errorf("expected pass: %s", results[0].Message)
	}
}

func TestMatchConstraints_NoTool_Fail(t *testing.T) {
	constraints := []spec.Constraint{{NoTool: "write_file"}}
	actual := []loader.ToolCall{{Name: "write_file", Arguments: map[string]interface{}{}}}
	results := MatchConstraints(constraints, nil, actual)
	if results[0].Passed {
		t.Error("expected fail: write_file was called")
	}
}

func TestMatchConstraints_MaxTools_Pass(t *testing.T) {
	constraints := []spec.Constraint{{MaxTools: 3}}
	actual := []loader.ToolCall{
		{Name: "a"}, {Name: "b"},
	}
	results := MatchConstraints(constraints, nil, actual)
	if !results[0].Passed {
		t.Errorf("expected pass: %s", results[0].Message)
	}
}

func TestMatchConstraints_MaxTools_Fail(t *testing.T) {
	constraints := []spec.Constraint{{MaxTools: 2}}
	actual := []loader.ToolCall{
		{Name: "a"}, {Name: "b"}, {Name: "c"},
	}
	results := MatchConstraints(constraints, nil, actual)
	if results[0].Passed {
		t.Error("expected fail: exceeded max tools")
	}
}

func TestMatchConstraints_MaxTools_Exact(t *testing.T) {
	constraints := []spec.Constraint{{MaxTools: 2}}
	actual := []loader.ToolCall{
		{Name: "a"}, {Name: "b"},
	}
	results := MatchConstraints(constraints, nil, actual)
	if !results[0].Passed {
		t.Errorf("expected pass at exact limit: %s", results[0].Message)
	}
}

func TestMatchConstraints_Ordered_Pass(t *testing.T) {
	constraints := []spec.Constraint{{Ordered: true}}
	expected := []spec.ToolExpect{{Name: "a"}, {Name: "b"}, {Name: "c"}}
	actual := []loader.ToolCall{
		{Name: "a"}, {Name: "b"}, {Name: "c"},
	}
	results := MatchConstraints(constraints, expected, actual)
	if !results[0].Passed {
		t.Errorf("expected pass: %s", results[0].Message)
	}
}

func TestMatchConstraints_Ordered_Fail(t *testing.T) {
	constraints := []spec.Constraint{{Ordered: true}}
	expected := []spec.ToolExpect{{Name: "a"}, {Name: "b"}}
	actual := []loader.ToolCall{
		{Name: "b"}, {Name: "a"},
	}
	results := MatchConstraints(constraints, expected, actual)
	if results[0].Passed {
		t.Error("expected fail: wrong order")
	}
}

func TestMatchConstraints_Ordered_WithGaps(t *testing.T) {
	constraints := []spec.Constraint{{Ordered: true}}
	expected := []spec.ToolExpect{{Name: "a"}, {Name: "c"}}
	actual := []loader.ToolCall{
		{Name: "a"}, {Name: "b"}, {Name: "c"},
	}
	results := MatchConstraints(constraints, expected, actual)
	if !results[0].Passed {
		t.Errorf("expected pass with gaps: %s", results[0].Message)
	}
}

func TestMatchConstraints_Ordered_MissingTool(t *testing.T) {
	constraints := []spec.Constraint{{Ordered: true}}
	expected := []spec.ToolExpect{{Name: "a"}, {Name: "missing"}}
	actual := []loader.ToolCall{{Name: "a"}}
	results := MatchConstraints(constraints, expected, actual)
	if results[0].Passed {
		t.Error("expected fail: missing tool in ordered check")
	}
}

func TestMatchConstraints_Multiple(t *testing.T) {
	constraints := []spec.Constraint{
		{NoTool: "delete"},
		{MaxTools: 5},
	}
	actual := []loader.ToolCall{{Name: "read"}, {Name: "write"}}
	results := MatchConstraints(constraints, nil, actual)
	if len(results) != 2 {
		t.Fatalf("results = %d, want 2", len(results))
	}
	if !results[0].Passed || !results[1].Passed {
		t.Error("both constraints should pass")
	}
}

func TestMatchConstraints_Empty(t *testing.T) {
	results := MatchConstraints(nil, nil, nil)
	if len(results) != 0 {
		t.Errorf("results = %d, want 0", len(results))
	}
}
