[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Tests](https://img.shields.io/badge/Tests-134-blue?style=for-the-badge)](https://github.com/JSLEEKR/agentspec)
[![Platform](https://img.shields.io/badge/Platform-Cross--platform-lightgrey?style=for-the-badge)]()

# agentspec

**YAML-driven behavioral testing for AI agents.**

Define what your agent *should do* in YAML. Run it against execution logs. Get pass/fail results. No LLM calls required.

---

## Why This Exists

Everyone builds AI agents. Nobody tests them systematically.

Developers test LLM output quality (promptfoo, DeepEval) or mock LLM calls (mocklm). But nobody tests agent **behavior** -- did the agent call the right tools? In the right order? Did it avoid calling dangerous tools?

`agentspec` fills this gap. It reads YAML behavior specs and checks them against execution logs. Deterministic. Fast. CI-friendly.

**The problem in one sentence:** You can test that GPT returns good text, but you can't test that your agent followed the correct workflow.

---

## How It Works

```
+------------------+     +------------------+     +------------------+
|   YAML Spec      | --> |   agentspec      | <-- |  Execution Log   |
|   (expected)     |     |   (engine)       |     |  (actual)        |
+------------------+     +------------------+     +------------------+
                               |
                               v
                    +------------------+
                    |   Pass / Fail    |
                    |   Report         |
                    +------------------+
```

1. Write a YAML spec defining expected agent behavior
2. Run your agent and capture the execution log (JSON)
3. `agentspec run specs/` compares spec vs log
4. Get structured pass/fail results (table or JSON)

---

## Installation

```bash
go install github.com/JSLEEKR/agentspec/cmd/agentspec@latest
```

Or build from source:

```bash
git clone https://github.com/JSLEEKR/agentspec.git
cd agentspec
go build -o agentspec ./cmd/agentspec
```

---

## Quick Start

### 1. Create example specs

```bash
agentspec init
```

This creates `specs/file-reader.yaml` and `specs/file-reader.log.json`.

### 2. Run specs

```bash
agentspec run specs/
```

Output:

```
agentspec — Agent Behavior Testing
===================================

specs/file-reader.yaml
  PASS Tool call: read_file (path=main.go)
  PASS Response contains "package main"
  PASS Constraint: no write_file calls
  PASS Constraint: max 3 tool calls (1 used)
  PASS Constraint: ordered execution

Results: 5 passed, 0 failed (1 specs)
```

### 3. Validate specs without running

```bash
agentspec validate specs/
```

---

## Writing Specs

### Basic Spec

```yaml
name: "File reader agent reads requested files"
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
```

### Matching Modes

agentspec supports 5 matching modes for tool arguments:

#### Exact Match (default)

```yaml
args:
  path: "main.go"  # must be exactly "main.go"
```

#### Contains

```yaml
args:
  query:
    contains: "weather"  # must contain "weather"
```

#### Regex

```yaml
args:
  query:
    regex: "weather.*tokyo"  # must match regex
```

#### Schema (Type Check)

```yaml
args:
  count:
    type: "number"  # must be a number
```

Supported types: `string`, `number`, `boolean`, `null`, `array`, `object`.

#### Any

```yaml
args:
  session_id:
    any: true  # matches any value
```

### Response Matching

```yaml
expect:
  response:
    contains: "package main"     # substring match
    exact: "Done."               # exact match
    regex: "\\d+ results found"  # regex match
```

### Constraints

```yaml
expect:
  constraints:
    - no_tool: "write_file"   # agent must NOT call write_file
    - max_tools: 3            # at most 3 tool calls total
    - ordered: true           # tools must be called in listed order
```

---

## Execution Logs

agentspec reads JSON execution logs that describe what the agent actually did:

```json
{
  "input": "Read the contents of main.go",
  "tool_calls": [
    {
      "name": "read_file",
      "arguments": {"path": "main.go"},
      "result": "package main\n\nfunc main() {}"
    }
  ],
  "response": "Here are the contents of main.go: package main..."
}
```

### Log Discovery

By default, agentspec looks for logs alongside specs:

```
specs/
  file-reader.yaml       <-- spec
  file-reader.log.json   <-- execution log (auto-discovered)
  search-agent.yaml
  search-agent.log.json
```

Or specify a log directory:

```bash
agentspec run specs/ --logs logs/
```

Or a single log for all specs:

```bash
agentspec run specs/ --logs execution.json
```

---

## CLI Reference

### `agentspec run <spec-path>`

Run specs against execution logs.

```bash
agentspec run specs/                    # run all specs in directory
agentspec run specs/file-reader.yaml    # run single spec
agentspec run specs/ --format json      # JSON output for CI
agentspec run specs/ --format table     # table output (default)
agentspec run specs/ --parallel 4       # parallel execution
agentspec run specs/ --logs logs/       # specify log directory
```

Exit code 1 if any spec fails (CI-friendly).

### `agentspec validate <spec-path>`

Validate spec syntax without running.

```bash
agentspec validate specs/
```

### `agentspec init`

Create example spec and execution log files.

```bash
agentspec init
```

### `agentspec version`

Show version.

```bash
agentspec version
```

---

## Output Formats

### Table (default)

```
agentspec — Agent Behavior Testing
===================================

specs/file-reader.yaml
  PASS Tool call: read_file (path=main.go)
  PASS Response contains "package main"
  PASS Constraint: no write_file calls
  PASS Constraint: max 3 tool calls (1 used)

specs/search-agent.yaml
  PASS Tool call: web_search (query matches /weather.*tokyo/)
  FAIL Tool call: summarize -- expected but not called
  PASS Constraint: ordered execution

Results: 5 passed, 1 failed (2 specs)
```

### JSON

```json
{
  "summary": {
    "total": 2,
    "passed": 1,
    "failed": 1,
    "checks": 6
  },
  "specs": [
    {
      "name": "File reader agent",
      "path": "specs/file-reader.yaml",
      "passed": true,
      "checks": [
        {"passed": true, "message": "Tool call: read_file (path=main.go)"}
      ]
    }
  ]
}
```

---

## Architecture

```
cmd/agentspec/main.go           -- CLI entry (cobra)
internal/spec/
  parser.go                     -- YAML spec parser
  types.go                      -- Spec data structures
  validator.go                  -- Spec syntax validation
internal/matcher/
  matcher.go                    -- Tool call matching engine
  pattern.go                    -- Exact/contains/regex/schema/any
internal/runner/
  runner.go                     -- Spec execution runner
  parallel.go                   -- Parallel execution
internal/reporter/
  table.go                      -- Table output
  json.go                       -- JSON output
  summary.go                    -- Pass/fail summary
internal/loader/
  loader.go                     -- Load execution logs (JSON)
```

---

## Use Cases

### CI Pipeline

```yaml
# .github/workflows/agent-test.yml
- name: Test agent behavior
  run: |
    go install github.com/JSLEEKR/agentspec/cmd/agentspec@latest
    agentspec run specs/ --format json > results.json
```

### MCP Server Testing

Test that your MCP server agent calls the right tools:

```yaml
name: "Code review agent uses diff tool"
input:
  message: "Review this pull request"
expect:
  tools:
    - name: git_diff
      args:
        ref:
          regex: "^(main|master)\\.\\.\\.HEAD$"
    - name: read_file
      args:
        path:
          any: true
  constraints:
    - no_tool: "git_push"
    - max_tools: 10
    - ordered: true
```

### Agent Workflow Validation

Ensure agents follow the correct sequence:

```yaml
name: "Research agent follows search-then-summarize pattern"
input:
  message: "Research quantum computing advances in 2026"
expect:
  tools:
    - name: web_search
      args:
        query:
          contains: "quantum"
    - name: summarize
  constraints:
    - ordered: true
    - no_tool: "write_file"
```

---

## Security

- No network calls -- reads local files only
- No LLM API calls -- decoupled by design
- YAML parsing with 1MB size limit per spec file
- No code execution from specs
- Strict JSON parsing mode rejects unknown fields

---

## Development

```bash
# Run tests
go test ./...

# Build
go build -o agentspec ./cmd/agentspec

# Run example
agentspec init
agentspec run specs/
```

---

## License

MIT License. See [LICENSE](LICENSE).
