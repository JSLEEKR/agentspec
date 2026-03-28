# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-03-28

### Added

- YAML-driven behavior spec parser with size limit enforcement (1MB)
- Spec validation with detailed error reporting
- Tool call matching engine with 5 matching modes:
  - `exact` -- literal string comparison
  - `contains` -- substring matching
  - `regex` -- regular expression matching
  - `schema` -- type-based validation (string, number, boolean, null, array, object)
  - `any` -- matches any value
- Behavioral constraints:
  - `no_tool` -- assert a tool was NOT called
  - `max_tools` -- limit total tool call count
  - `ordered` -- assert tools called in specified order
- Response matching (contains, exact, regex)
- Execution log loader (JSON format) with strict and flexible modes
- Sequential and parallel spec execution
- Table and JSON output formats
- CLI commands: `run`, `validate`, `init`, `version`
- Auto-discovery of execution logs alongside spec files (`*.log.json`)
- 134 tests across all packages
