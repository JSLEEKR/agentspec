---
name: test-agent
description: Run behavioral specs against AI agent execution logs. Use when testing agent tool-call patterns, verifying agent compliance with specs, or running CI checks on agent behavior.
---

# Test Agent Behavior

Run YAML behavioral specs against agent execution logs.

## Usage
```bash
# Run all specs
agentspec run specs/

# Run single spec
agentspec run specs/file-reader.yaml

# JSON output for CI
agentspec run specs/ --format json

# Validate spec syntax
agentspec validate specs/

# Create example specs
agentspec init
```

## Spec Format
```yaml
name: "Agent reads requested files"
input:
  message: "Read main.go"
expect:
  tools:
    - name: read_file
      args:
        path: "main.go"
  constraints:
    - no_tool: "write_file"
    - max_tools: 3
```
