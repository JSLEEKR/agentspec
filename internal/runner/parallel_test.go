package runner

import (
	"testing"

	"github.com/JSLEEKR/agentspec/internal/loader"
	"github.com/JSLEEKR/agentspec/internal/spec"
)

func TestRunParallel_AllPass(t *testing.T) {
	n := 10
	specs := make([]*spec.Spec, n)
	paths := make([]string, n)
	logs := make([]*loader.ExecutionLog, n)
	for i := 0; i < n; i++ {
		specs[i] = makeSpec("test", []spec.ToolExpect{{Name: "tool"}}, nil, nil)
		paths[i] = "test.yaml"
		logs[i] = makeLog([]loader.ToolCall{{Name: "tool", Arguments: map[string]interface{}{}}}, "")
	}

	result, err := RunParallel(specs, paths, logs, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalPassed != n {
		t.Errorf("passed = %d, want %d", result.TotalPassed, n)
	}
	if result.TotalFailed != 0 {
		t.Errorf("failed = %d, want 0", result.TotalFailed)
	}
}

func TestRunParallel_Mixed(t *testing.T) {
	specs := []*spec.Spec{
		makeSpec("pass", []spec.ToolExpect{{Name: "t"}}, nil, nil),
		makeSpec("fail", []spec.ToolExpect{{Name: "missing"}}, nil, nil),
		makeSpec("pass2", []spec.ToolExpect{{Name: "t"}}, nil, nil),
	}
	paths := []string{"a.yaml", "b.yaml", "c.yaml"}
	logs := []*loader.ExecutionLog{
		makeLog([]loader.ToolCall{{Name: "t", Arguments: map[string]interface{}{}}}, ""),
		makeLog([]loader.ToolCall{{Name: "other", Arguments: map[string]interface{}{}}}, ""),
		makeLog([]loader.ToolCall{{Name: "t", Arguments: map[string]interface{}{}}}, ""),
	}

	result, err := RunParallel(specs, paths, logs, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalPassed != 2 {
		t.Errorf("passed = %d, want 2", result.TotalPassed)
	}
	if result.TotalFailed != 1 {
		t.Errorf("failed = %d, want 1", result.TotalFailed)
	}
}

func TestRunParallel_SingleWorker(t *testing.T) {
	specs := []*spec.Spec{
		makeSpec("test", []spec.ToolExpect{{Name: "t"}}, nil, nil),
	}
	logs := []*loader.ExecutionLog{
		makeLog([]loader.ToolCall{{Name: "t", Arguments: map[string]interface{}{}}}, ""),
	}

	result, err := RunParallel(specs, nil, logs, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalPassed != 1 {
		t.Errorf("passed = %d, want 1", result.TotalPassed)
	}
}

func TestRunParallel_ZeroWorkers(t *testing.T) {
	specs := []*spec.Spec{
		makeSpec("test", []spec.ToolExpect{{Name: "t"}}, nil, nil),
	}
	logs := []*loader.ExecutionLog{
		makeLog([]loader.ToolCall{{Name: "t", Arguments: map[string]interface{}{}}}, ""),
	}

	// 0 workers should default to 1
	result, err := RunParallel(specs, nil, logs, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalPassed != 1 {
		t.Errorf("passed = %d, want 1", result.TotalPassed)
	}
}

func TestRunParallel_MoreWorkersThanSpecs(t *testing.T) {
	specs := []*spec.Spec{
		makeSpec("test", []spec.ToolExpect{{Name: "t"}}, nil, nil),
	}
	logs := []*loader.ExecutionLog{
		makeLog([]loader.ToolCall{{Name: "t", Arguments: map[string]interface{}{}}}, ""),
	}

	result, err := RunParallel(specs, nil, logs, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalPassed != 1 {
		t.Errorf("passed = %d, want 1", result.TotalPassed)
	}
}

func TestRunParallel_PreservesOrder(t *testing.T) {
	n := 5
	specs := make([]*spec.Spec, n)
	logs := make([]*loader.ExecutionLog, n)
	for i := 0; i < n; i++ {
		name := "pass"
		toolName := "t"
		if i%2 == 1 {
			name = "fail"
			toolName = "missing"
		}
		specs[i] = makeSpec(name, []spec.ToolExpect{{Name: toolName}}, nil, nil)
		logs[i] = makeLog([]loader.ToolCall{{Name: "t", Arguments: map[string]interface{}{}}}, "")
	}

	result, err := RunParallel(specs, nil, logs, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, sr := range result.Specs {
		if i%2 == 0 {
			if !sr.Passed {
				t.Errorf("spec[%d] should pass", i)
			}
		} else {
			if sr.Passed {
				t.Errorf("spec[%d] should fail", i)
			}
		}
	}
}
