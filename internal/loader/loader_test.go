package loader

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad_Valid(t *testing.T) {
	json := `{
		"input": "do something",
		"tool_calls": [
			{
				"name": "read_file",
				"arguments": {"path": "main.go"},
				"result": "package main"
			}
		],
		"response": "Here is the file"
	}`
	log, err := Load(strings.NewReader(json))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if log.Input != "do something" {
		t.Errorf("input = %q", log.Input)
	}
	if len(log.ToolCalls) != 1 {
		t.Fatalf("tool_calls = %d, want 1", len(log.ToolCalls))
	}
	if log.ToolCalls[0].Name != "read_file" {
		t.Errorf("tool name = %q", log.ToolCalls[0].Name)
	}
	if log.Response != "Here is the file" {
		t.Errorf("response = %q", log.Response)
	}
}

func TestLoad_MultipleToolCalls(t *testing.T) {
	json := `{
		"input": "test",
		"tool_calls": [
			{"name": "search", "arguments": {"q": "hello"}, "result": "found"},
			{"name": "write", "arguments": {"path": "out.txt"}, "result": "ok"}
		],
		"response": "done"
	}`
	log, err := Load(strings.NewReader(json))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(log.ToolCalls) != 2 {
		t.Errorf("tool_calls = %d, want 2", len(log.ToolCalls))
	}
}

func TestLoad_EmptyToolCalls(t *testing.T) {
	json := `{"input": "test", "tool_calls": [], "response": "ok"}`
	log, err := Load(strings.NewReader(json))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(log.ToolCalls) != 0 {
		t.Errorf("tool_calls = %d, want 0", len(log.ToolCalls))
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	_, err := Load(strings.NewReader("{invalid"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoad_NoToolCalls(t *testing.T) {
	json := `{"input": "test", "response": "ok"}`
	log, err := Load(strings.NewReader(json))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if log.ToolCalls != nil {
		t.Errorf("tool_calls should be nil, got %v", log.ToolCalls)
	}
}

func TestLoadFlexible_ExtraFields(t *testing.T) {
	json := `{
		"input": "test",
		"tool_calls": [],
		"response": "ok",
		"extra_field": "should be ignored"
	}`
	log, err := LoadFlexible(strings.NewReader(json))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if log.Input != "test" {
		t.Errorf("input = %q", log.Input)
	}
}

func TestLoad_StrictRejectsExtraFields(t *testing.T) {
	json := `{
		"input": "test",
		"tool_calls": [],
		"response": "ok",
		"extra": "bad"
	}`
	_, err := Load(strings.NewReader(json))
	if err == nil {
		t.Error("expected error for extra fields in strict mode")
	}
}

func TestLoadFile_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.json")
	content := `{"input": "test", "tool_calls": [], "response": "ok"}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	log, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if log.Input != "test" {
		t.Errorf("input = %q", log.Input)
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	_, err := LoadFile("/nonexistent/file.json")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestLoad_ToolArgTypes(t *testing.T) {
	json := `{
		"input": "test",
		"tool_calls": [
			{
				"name": "tool",
				"arguments": {
					"str": "hello",
					"num": 42,
					"bool": true,
					"null_val": null
				},
				"result": "ok"
			}
		],
		"response": "ok"
	}`
	log, err := Load(strings.NewReader(json))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	args := log.ToolCalls[0].Arguments
	if _, ok := args["str"].(string); !ok {
		t.Error("str should be string")
	}
	if _, ok := args["num"].(float64); !ok {
		t.Error("num should be float64")
	}
	if _, ok := args["bool"].(bool); !ok {
		t.Error("bool should be bool")
	}
	if args["null_val"] != nil {
		t.Error("null_val should be nil")
	}
}

func TestLoad_NestedResult(t *testing.T) {
	json := `{
		"input": "test",
		"tool_calls": [
			{
				"name": "tool",
				"arguments": {},
				"result": {"key": "value"}
			}
		],
		"response": "ok"
	}`
	log, err := Load(strings.NewReader(json))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, ok := log.ToolCalls[0].Result.(map[string]interface{})
	if !ok {
		t.Fatal("result should be a map")
	}
	if result["key"] != "value" {
		t.Errorf("result.key = %v", result["key"])
	}
}
