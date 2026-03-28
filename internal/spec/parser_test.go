package spec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParse_ValidSpec(t *testing.T) {
	yaml := `
name: "test spec"
input:
  message: "do something"
expect:
  tools:
    - name: read_file
      args:
        path: "main.go"
`
	s, err := Parse(strings.NewReader(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Name != "test spec" {
		t.Errorf("name = %q, want %q", s.Name, "test spec")
	}
	if s.Input.Message != "do something" {
		t.Errorf("input.message = %q, want %q", s.Input.Message, "do something")
	}
	if len(s.Expect.Tools) != 1 {
		t.Fatalf("tools count = %d, want 1", len(s.Expect.Tools))
	}
	if s.Expect.Tools[0].Name != "read_file" {
		t.Errorf("tool name = %q, want %q", s.Expect.Tools[0].Name, "read_file")
	}
}

func TestParse_WithResponse(t *testing.T) {
	yaml := `
name: "response test"
input:
  message: "test"
expect:
  tools:
    - name: foo
  response:
    contains: "hello"
`
	s, err := Parse(strings.NewReader(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Expect.Response == nil {
		t.Fatal("response is nil")
	}
	if s.Expect.Response.Contains != "hello" {
		t.Errorf("response.contains = %q, want %q", s.Expect.Response.Contains, "hello")
	}
}

func TestParse_WithConstraints(t *testing.T) {
	yaml := `
name: "constraint test"
input:
  message: "test"
expect:
  tools:
    - name: foo
  constraints:
    - no_tool: "bar"
    - max_tools: 5
    - ordered: true
`
	s, err := Parse(strings.NewReader(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Expect.Constraints) != 3 {
		t.Fatalf("constraints count = %d, want 3", len(s.Expect.Constraints))
	}
	if s.Expect.Constraints[0].NoTool != "bar" {
		t.Errorf("constraint[0].no_tool = %q, want %q", s.Expect.Constraints[0].NoTool, "bar")
	}
	if s.Expect.Constraints[1].MaxTools != 5 {
		t.Errorf("constraint[1].max_tools = %d, want 5", s.Expect.Constraints[1].MaxTools)
	}
	if !s.Expect.Constraints[2].Ordered {
		t.Error("constraint[2].ordered = false, want true")
	}
}

func TestParse_InvalidYAML(t *testing.T) {
	_, err := Parse(strings.NewReader("{{invalid"))
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestParse_EmptySpec(t *testing.T) {
	s, err := Parse(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Name != "" {
		t.Errorf("name = %q, want empty", s.Name)
	}
}

func TestParse_MultipleTools(t *testing.T) {
	yaml := `
name: "multi"
input:
  message: "test"
expect:
  tools:
    - name: search
      args:
        query: "hello"
    - name: summarize
    - name: write_file
      args:
        path: "out.txt"
        content: "result"
`
	s, err := Parse(strings.NewReader(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Expect.Tools) != 3 {
		t.Fatalf("tools count = %d, want 3", len(s.Expect.Tools))
	}
}

func TestParse_MatcherArgs(t *testing.T) {
	yaml := `
name: "matchers"
input:
  message: "test"
expect:
  tools:
    - name: search
      args:
        query:
          contains: "weather"
        location:
          regex: ".*tokyo.*"
        anything:
          any: true
        count:
          type: "number"
`
	s, err := Parse(strings.NewReader(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	args := s.Expect.Tools[0].Args
	if len(args) != 4 {
		t.Fatalf("args count = %d, want 4", len(args))
	}
}

func TestParse_ResponseRegex(t *testing.T) {
	yaml := `
name: "regex response"
input:
  message: "test"
expect:
  response:
    regex: "\\d+ results"
`
	s, err := Parse(strings.NewReader(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Expect.Response.Regex != `\d+ results` {
		t.Errorf("response.regex = %q", s.Expect.Response.Regex)
	}
}

func TestParse_ResponseExact(t *testing.T) {
	yaml := `
name: "exact response"
input:
  message: "test"
expect:
  response:
    exact: "Done."
`
	s, err := Parse(strings.NewReader(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Expect.Response.Exact != "Done." {
		t.Errorf("response.exact = %q", s.Expect.Response.Exact)
	}
}

func TestParseFile_NotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/spec.yaml")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestParseFile_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.yaml")
	content := `
name: "file test"
input:
  message: "hello"
expect:
  tools:
    - name: greet
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	s, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Name != "file test" {
		t.Errorf("name = %q", s.Name)
	}
}

func TestParseDir_Empty(t *testing.T) {
	dir := t.TempDir()
	specs, paths, err := ParseDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(specs) != 0 {
		t.Errorf("specs = %d, want 0", len(specs))
	}
	if len(paths) != 0 {
		t.Errorf("paths = %d, want 0", len(paths))
	}
}

func TestParseDir_WithSpecs(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a.yaml", "b.yml", "c.txt"} {
		content := `
name: "` + name + `"
input:
  message: "test"
expect:
  tools:
    - name: foo
`
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	specs, paths, err := ParseDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(specs) != 2 {
		t.Errorf("specs = %d, want 2 (yaml + yml only)", len(specs))
	}
	if len(paths) != 2 {
		t.Errorf("paths = %d, want 2", len(paths))
	}
}

func TestParseDir_NotFound(t *testing.T) {
	_, _, err := ParseDir("/nonexistent/dir")
	if err == nil {
		t.Error("expected error for nonexistent dir")
	}
}

func TestParse_ToolWithNoArgs(t *testing.T) {
	yaml := `
name: "no args"
input:
  message: "test"
expect:
  tools:
    - name: list_files
`
	s, err := Parse(strings.NewReader(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Expect.Tools[0].Args) != 0 {
		t.Errorf("args should be empty, got %v", s.Expect.Tools[0].Args)
	}
}

func TestParseFile_TooLarge(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "huge.yaml")
	// Create a file just over 1MB
	data := make([]byte, maxSpecSize+100)
	for i := range data {
		data[i] = 'a'
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := ParseFile(path)
	if err == nil {
		t.Error("expected error for oversized file")
	}
}

func TestParse_OnlyConstraints(t *testing.T) {
	yaml := `
name: "constraints only"
input:
  message: "test"
expect:
  constraints:
    - max_tools: 2
`
	s, err := Parse(strings.NewReader(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Expect.Constraints) != 1 {
		t.Errorf("constraints = %d, want 1", len(s.Expect.Constraints))
	}
}

func TestParse_OnlyResponse(t *testing.T) {
	yaml := `
name: "response only"
input:
  message: "test"
expect:
  response:
    contains: "answer"
`
	s, err := Parse(strings.NewReader(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Expect.Response == nil {
		t.Fatal("response should not be nil")
	}
}
