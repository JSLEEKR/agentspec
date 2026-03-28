package loader

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// ExecutionLog represents an agent's execution log.
type ExecutionLog struct {
	Input     string     `json:"input"`
	ToolCalls []ToolCall `json:"tool_calls"`
	Response  string     `json:"response"`
}

// ToolCall represents a single tool invocation by the agent.
type ToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
	Result    interface{}            `json:"result"`
}

// LoadFile loads an execution log from a JSON file.
func LoadFile(path string) (*ExecutionLog, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open log %s: %w", path, err)
	}
	defer f.Close()
	return Load(f)
}

// Load parses an execution log from a reader, allowing unknown fields
// for compatibility with real-world logs that may contain extra data.
func Load(r io.Reader) (*ExecutionLog, error) {
	var log ExecutionLog
	dec := json.NewDecoder(r)
	if err := dec.Decode(&log); err != nil {
		return nil, fmt.Errorf("parse execution log: %w", err)
	}
	return &log, nil
}

// LoadStrict parses an execution log, rejecting unknown fields.
func LoadStrict(r io.Reader) (*ExecutionLog, error) {
	var log ExecutionLog
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&log); err != nil {
		return nil, fmt.Errorf("parse execution log: %w", err)
	}
	return &log, nil
}
