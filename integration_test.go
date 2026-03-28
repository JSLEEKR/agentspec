package agentspec_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JSLEEKR/agentspec/internal/loader"
	"github.com/JSLEEKR/agentspec/internal/reporter"
	"github.com/JSLEEKR/agentspec/internal/runner"
	"github.com/JSLEEKR/agentspec/internal/spec"
)

func TestIntegration_FullPipeline_Pass(t *testing.T) {
	dir := t.TempDir()

	// Write spec
	specContent := `
name: "File reader reads main.go"
input:
  message: "Read the contents of main.go"
expect:
  tools:
    - name: read_file
      args:
        path: "main.go"
  response:
    contains: "package main"
  constraints:
    - no_tool: "write_file"
    - max_tools: 3
    - ordered: true
`
	specPath := filepath.Join(dir, "test.yaml")
	os.WriteFile(specPath, []byte(specContent), 0o644)

	// Write log
	logContent := `{
  "input": "Read the contents of main.go",
  "tool_calls": [
    {
      "name": "read_file",
      "arguments": {"path": "main.go"},
      "result": "package main\n\nfunc main() {}"
    }
  ],
  "response": "Here are the contents of main.go: package main..."
}`
	logPath := filepath.Join(dir, "test.log.json")
	os.WriteFile(logPath, []byte(logContent), 0o644)

	// Parse spec
	s, err := spec.ParseFile(specPath)
	if err != nil {
		t.Fatalf("parse spec: %v", err)
	}

	// Validate spec
	if err := spec.Validate(s); err != nil {
		t.Fatalf("validate spec: %v", err)
	}

	// Load log
	log, err := loader.LoadFile(logPath)
	if err != nil {
		t.Fatalf("load log: %v", err)
	}

	// Run
	result := runner.RunSingle(s, log)

	if result.TotalPassed != 1 {
		t.Errorf("passed = %d, want 1", result.TotalPassed)
	}
	if result.TotalFailed != 0 {
		t.Errorf("failed = %d, want 0", result.TotalFailed)
	}

	// Verify table output
	output := reporter.FormatTableString(result)
	if !strings.Contains(output, "PASS") {
		t.Errorf("table output should contain PASS:\n%s", output)
	}
}

func TestIntegration_FullPipeline_Fail(t *testing.T) {
	dir := t.TempDir()

	specContent := `
name: "Search agent searches and summarizes"
input:
  message: "Search for weather"
expect:
  tools:
    - name: web_search
      args:
        query:
          contains: "weather"
    - name: summarize
  constraints:
    - ordered: true
`
	specPath := filepath.Join(dir, "search.yaml")
	os.WriteFile(specPath, []byte(specContent), 0o644)

	logContent := `{
  "input": "Search for weather",
  "tool_calls": [
    {
      "name": "web_search",
      "arguments": {"query": "weather in tokyo"},
      "result": "sunny 25C"
    }
  ],
  "response": "The weather in Tokyo is sunny."
}`
	logPath := filepath.Join(dir, "search.log.json")
	os.WriteFile(logPath, []byte(logContent), 0o644)

	s, _ := spec.ParseFile(specPath)
	spec.Validate(s)
	log, _ := loader.LoadFile(logPath)

	result := runner.RunSingle(s, log)

	if result.TotalFailed != 1 {
		t.Errorf("failed = %d, want 1 (summarize not called)", result.TotalFailed)
	}
}

func TestIntegration_Directory(t *testing.T) {
	dir := t.TempDir()

	// Two specs
	for _, name := range []string{"a", "b"} {
		specContent := `
name: "` + name + `"
input:
  message: "test"
expect:
  tools:
    - name: tool_` + name + `
`
		logContent := `{
  "input": "test",
  "tool_calls": [{"name": "tool_` + name + `", "arguments": {}, "result": "ok"}],
  "response": "done"
}`
		os.WriteFile(filepath.Join(dir, name+".yaml"), []byte(specContent), 0o644)
		os.WriteFile(filepath.Join(dir, name+".log.json"), []byte(logContent), 0o644)
	}

	specs, paths, err := spec.ParseDir(dir)
	if err != nil {
		t.Fatalf("parse dir: %v", err)
	}

	var logs []*loader.ExecutionLog
	for _, p := range paths {
		ext := filepath.Ext(p)
		logFile := strings.TrimSuffix(p, ext) + ".log.json"
		log, err := loader.LoadFile(logFile)
		if err != nil {
			t.Fatalf("load log: %v", err)
		}
		logs = append(logs, log)
	}

	result, err := runner.RunAll(specs, paths, logs)
	if err != nil {
		t.Fatalf("run all: %v", err)
	}

	if result.TotalPassed != 2 {
		t.Errorf("passed = %d, want 2", result.TotalPassed)
	}
}

func TestIntegration_Parallel(t *testing.T) {
	dir := t.TempDir()

	n := 5
	specs := make([]*spec.Spec, n)
	paths := make([]string, n)
	logs := make([]*loader.ExecutionLog, n)

	for i := 0; i < n; i++ {
		specContent := `
name: "parallel test"
input:
  message: "test"
expect:
  tools:
    - name: tool
`
		logContent := `{"input": "test", "tool_calls": [{"name": "tool", "arguments": {}, "result": "ok"}], "response": "ok"}`
		sp := filepath.Join(dir, "spec"+string(rune('a'+i))+".yaml")
		lp := filepath.Join(dir, "spec"+string(rune('a'+i))+".log.json")
		os.WriteFile(sp, []byte(specContent), 0o644)
		os.WriteFile(lp, []byte(logContent), 0o644)

		s, _ := spec.ParseFile(sp)
		l, _ := loader.LoadFile(lp)
		specs[i] = s
		paths[i] = sp
		logs[i] = l
	}

	result := runner.RunParallel(specs, paths, logs, 3)
	if result.TotalPassed != n {
		t.Errorf("passed = %d, want %d", result.TotalPassed, n)
	}
}

func TestIntegration_JSONOutput(t *testing.T) {
	s := &spec.Spec{
		Name:  "json test",
		Input: spec.Input{Message: "test"},
		Expect: spec.Expect{
			Tools: []spec.ToolExpect{{Name: "tool"}},
		},
	}
	log := &loader.ExecutionLog{
		Input:     "test",
		ToolCalls: []loader.ToolCall{{Name: "tool", Arguments: map[string]interface{}{}}},
		Response:  "ok",
	}

	result := runner.RunSingle(s, log)

	var sb strings.Builder
	if err := reporter.FormatJSON(&sb, result); err != nil {
		t.Fatalf("format json: %v", err)
	}

	output := sb.String()
	if !strings.Contains(output, `"passed": true`) {
		t.Errorf("JSON should contain passed: true:\n%s", output)
	}
}

func TestIntegration_ConstraintViolation(t *testing.T) {
	s := &spec.Spec{
		Name:  "constraint test",
		Input: spec.Input{Message: "test"},
		Expect: spec.Expect{
			Constraints: []spec.Constraint{
				{NoTool: "dangerous_tool"},
				{MaxTools: 2},
			},
		},
	}
	log := &loader.ExecutionLog{
		Input: "test",
		ToolCalls: []loader.ToolCall{
			{Name: "safe_tool", Arguments: map[string]interface{}{}},
			{Name: "dangerous_tool", Arguments: map[string]interface{}{}},
			{Name: "another_tool", Arguments: map[string]interface{}{}},
		},
		Response: "done",
	}

	result := runner.RunSingle(s, log)
	if result.TotalFailed != 1 {
		t.Errorf("should fail: dangerous_tool called AND max_tools exceeded")
	}

	// Check that both constraints failed
	failCount := 0
	for _, c := range result.Specs[0].Checks {
		if !c.Passed {
			failCount++
		}
	}
	if failCount != 2 {
		t.Errorf("failed checks = %d, want 2", failCount)
	}
}

func TestIntegration_RegexMatching(t *testing.T) {
	s := &spec.Spec{
		Name:  "regex test",
		Input: spec.Input{Message: "test"},
		Expect: spec.Expect{
			Tools: []spec.ToolExpect{{
				Name: "search",
				Args: map[string]interface{}{
					"query": map[string]interface{}{"regex": `weather.*tokyo`},
				},
			}},
			Response: &spec.Response{
				Regex: `\d+.*degrees`,
			},
		},
	}
	log := &loader.ExecutionLog{
		Input: "test",
		ToolCalls: []loader.ToolCall{{
			Name:      "search",
			Arguments: map[string]interface{}{"query": "weather forecast for tokyo"},
		}},
		Response: "It is 25 degrees in Tokyo",
	}

	result := runner.RunSingle(s, log)
	if result.TotalFailed != 0 {
		t.Error("regex matching should pass")
		for _, c := range result.Specs[0].Checks {
			t.Logf("  %v: %s", c.Passed, c.Message)
		}
	}
}
