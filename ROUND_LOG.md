# ROUND_LOG

## Round 24: agentspec

- **Category:** Testing / AI DevTools
- **Language:** Go
- **Date:** 2026-03-28
- **Tests:** 134
- **Status:** Build complete

### What It Does

YAML-driven behavioral testing framework for AI agents. Verifies tool calls, sequences, and constraints against execution logs. Decoupled from LLMs -- reads JSON execution logs, not API responses.

### Architecture

- `internal/spec/` -- YAML parser, types, validator
- `internal/matcher/` -- 5 matching modes (exact, contains, regex, schema, any)
- `internal/runner/` -- Sequential and parallel execution
- `internal/reporter/` -- Table and JSON output
- `internal/loader/` -- Execution log loader
- `cmd/agentspec/` -- CLI (cobra)

### Key Decisions

1. **Decoupled from LLMs** -- Tests run against execution logs, not live agents. This makes tests deterministic and fast.
2. **YAML specs** -- Familiar format for defining expected behavior. Supports nested matchers.
3. **5 matching modes** -- From strict (exact) to permissive (any), covering real-world needs.
4. **Constraint system** -- `no_tool`, `max_tools`, `ordered` let you define what agents should NOT do.
5. **Auto log discovery** -- `spec.yaml` automatically pairs with `spec.log.json`.
