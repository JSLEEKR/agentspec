package runner

import (
	"testing"

	"github.com/JSLEEKR/agentspec/internal/loader"
	"github.com/JSLEEKR/agentspec/internal/spec"
)

func makeSpec(name string, tools []spec.ToolExpect, resp *spec.Response, constraints []spec.Constraint) *spec.Spec {
	return &spec.Spec{
		Name:  name,
		Input: spec.Input{Message: "test"},
		Expect: spec.Expect{
			Tools:       tools,
			Response:    resp,
			Constraints: constraints,
		},
	}
}

func makeLog(toolCalls []loader.ToolCall, response string) *loader.ExecutionLog {
	return &loader.ExecutionLog{
		Input:     "test",
		ToolCalls: toolCalls,
		Response:  response,
	}
}

func TestRun_AllPass(t *testing.T) {
	s := makeSpec("pass test",
		[]spec.ToolExpect{{Name: "read_file", Args: map[string]interface{}{"path": "main.go"}}},
		&spec.Response{Contains: "package main"},
		[]spec.Constraint{{NoTool: "write_file"}, {MaxTools: 3}},
	)
	log := makeLog(
		[]loader.ToolCall{{Name: "read_file", Arguments: map[string]interface{}{"path": "main.go"}}},
		"Here is main.go: package main...",
	)
	result := Run(s, log)
	if !result.Passed {
		t.Error("expected pass")
		for _, c := range result.Checks {
			t.Logf("  %v: %s", c.Passed, c.Message)
		}
	}
}

func TestRun_ToolFail(t *testing.T) {
	s := makeSpec("fail test",
		[]spec.ToolExpect{{Name: "missing_tool"}},
		nil, nil,
	)
	log := makeLog([]loader.ToolCall{{Name: "other_tool", Arguments: map[string]interface{}{}}}, "")
	result := Run(s, log)
	if result.Passed {
		t.Error("expected fail")
	}
}

func TestRun_ResponseFail(t *testing.T) {
	s := makeSpec("resp fail",
		nil,
		&spec.Response{Contains: "not here"},
		nil,
	)
	log := makeLog(nil, "something else")
	result := Run(s, log)
	if result.Passed {
		t.Error("expected fail")
	}
}

func TestRun_ConstraintFail(t *testing.T) {
	s := makeSpec("constraint fail",
		nil, nil,
		[]spec.Constraint{{NoTool: "bad_tool"}},
	)
	log := makeLog(
		[]loader.ToolCall{{Name: "bad_tool", Arguments: map[string]interface{}{}}},
		"",
	)
	result := Run(s, log)
	if result.Passed {
		t.Error("expected fail")
	}
}

func TestRunAll_Mixed(t *testing.T) {
	specs := []*spec.Spec{
		makeSpec("pass", []spec.ToolExpect{{Name: "tool"}}, nil, nil),
		makeSpec("fail", []spec.ToolExpect{{Name: "missing"}}, nil, nil),
	}
	paths := []string{"pass.yaml", "fail.yaml"}
	logs := []*loader.ExecutionLog{
		makeLog([]loader.ToolCall{{Name: "tool", Arguments: map[string]interface{}{}}}, ""),
		makeLog([]loader.ToolCall{{Name: "other", Arguments: map[string]interface{}{}}}, ""),
	}
	result, err := RunAll(specs, paths, logs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalPassed != 1 {
		t.Errorf("passed = %d, want 1", result.TotalPassed)
	}
	if result.TotalFailed != 1 {
		t.Errorf("failed = %d, want 1", result.TotalFailed)
	}
}

func TestRunAll_MismatchedCounts(t *testing.T) {
	specs := []*spec.Spec{makeSpec("test", nil, nil, nil)}
	logs := []*loader.ExecutionLog{makeLog(nil, ""), makeLog(nil, "")}
	_, err := RunAll(specs, nil, logs)
	if err == nil {
		t.Error("expected error for mismatched counts")
	}
}

func TestRunSingle(t *testing.T) {
	s := makeSpec("single", []spec.ToolExpect{{Name: "tool"}}, nil, nil)
	log := makeLog([]loader.ToolCall{{Name: "tool", Arguments: map[string]interface{}{}}}, "")
	result := RunSingle(s, log)
	if result.TotalPassed != 1 {
		t.Errorf("passed = %d, want 1", result.TotalPassed)
	}
	if len(result.Specs) != 1 {
		t.Errorf("specs = %d, want 1", len(result.Specs))
	}
}

func TestRunSingle_Fail(t *testing.T) {
	s := makeSpec("single fail", []spec.ToolExpect{{Name: "missing"}}, nil, nil)
	log := makeLog(nil, "")
	result := RunSingle(s, log)
	if result.TotalFailed != 1 {
		t.Errorf("failed = %d, want 1", result.TotalFailed)
	}
}

func TestRun_EmptySpec(t *testing.T) {
	s := makeSpec("empty", nil, nil, nil)
	log := makeLog(nil, "")
	result := Run(s, log)
	if !result.Passed {
		t.Error("empty spec should pass")
	}
	if len(result.Checks) != 0 {
		t.Errorf("checks = %d, want 0", len(result.Checks))
	}
}

func TestRun_SpecName(t *testing.T) {
	s := makeSpec("my test", []spec.ToolExpect{{Name: "t"}}, nil, nil)
	log := makeLog([]loader.ToolCall{{Name: "t", Arguments: map[string]interface{}{}}}, "")
	result := Run(s, log)
	if result.SpecName != "my test" {
		t.Errorf("specName = %q, want %q", result.SpecName, "my test")
	}
}

func TestRunAll_Paths(t *testing.T) {
	specs := []*spec.Spec{makeSpec("test", []spec.ToolExpect{{Name: "t"}}, nil, nil)}
	paths := []string{"specs/test.yaml"}
	logs := []*loader.ExecutionLog{makeLog([]loader.ToolCall{{Name: "t", Arguments: map[string]interface{}{}}}, "")}
	result, err := RunAll(specs, paths, logs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Specs[0].SpecPath != "specs/test.yaml" {
		t.Errorf("specPath = %q", result.Specs[0].SpecPath)
	}
}
