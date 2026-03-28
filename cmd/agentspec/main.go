package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/JSLEEKR/agentspec/internal/loader"
	"github.com/JSLEEKR/agentspec/internal/reporter"
	"github.com/JSLEEKR/agentspec/internal/runner"
	"github.com/JSLEEKR/agentspec/internal/spec"
	"github.com/spf13/cobra"
)

var version = "0.1.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "agentspec",
		Short: "Agent behavioral testing -- YAML specs for tool calls, sequences, constraints",
		Long:  "agentspec verifies AI agent behavior by running YAML specs against execution logs.",
	}

	// run command
	var formatFlag string
	var parallelFlag int
	var logDir string

	runCmd := &cobra.Command{
		Use:   "run <spec-path>",
		Short: "Run specs against execution logs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			specPath := args[0]
			return runSpecs(specPath, logDir, formatFlag, parallelFlag)
		},
	}
	runCmd.Flags().StringVar(&formatFlag, "format", "table", "Output format: table, json")
	runCmd.Flags().IntVar(&parallelFlag, "parallel", 1, "Number of parallel workers")
	runCmd.Flags().StringVar(&logDir, "logs", "", "Directory or file containing execution logs")

	// validate command
	validateCmd := &cobra.Command{
		Use:   "validate <spec-path>",
		Short: "Validate spec syntax without running",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return validateSpecs(args[0])
		},
	}

	// init command
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Create example spec files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return initSpecs()
		},
	}

	// version command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("agentspec %s\n", version)
		},
	}

	rootCmd.AddCommand(runCmd, validateCmd, initCmd, versionCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runSpecs(specPath, logPath, format string, parallel int) error {
	// Determine if specPath is file or directory
	info, err := os.Stat(specPath)
	if err != nil {
		return fmt.Errorf("cannot access %s: %w", specPath, err)
	}

	var specs []*spec.Spec
	var paths []string

	if info.IsDir() {
		specs, paths, err = spec.ParseDir(specPath)
		if err != nil {
			return err
		}
	} else {
		s, err := spec.ParseFile(specPath)
		if err != nil {
			return err
		}
		specs = []*spec.Spec{s}
		paths = []string{specPath}
	}

	if len(specs) == 0 {
		fmt.Println("No specs found.")
		return nil
	}

	// Validate all specs
	for i, s := range specs {
		if err := spec.Validate(s); err != nil {
			return fmt.Errorf("invalid spec %s: %w", paths[i], err)
		}
	}

	// Load execution logs
	logs, err := loadLogs(logPath, paths)
	if err != nil {
		return err
	}

	// Run specs
	var result *runner.RunResult
	if parallel > 1 {
		result = runner.RunParallel(specs, paths, logs, parallel)
	} else {
		result, err = runner.RunAll(specs, paths, logs)
		if err != nil {
			return err
		}
	}

	// Output results
	switch strings.ToLower(format) {
	case "json":
		if err := reporter.FormatJSON(os.Stdout, result); err != nil {
			return err
		}
	default:
		reporter.FormatTable(os.Stdout, result)
	}

	if result.TotalFailed > 0 {
		os.Exit(1)
	}
	return nil
}

func loadLogs(logPath string, specPaths []string) ([]*loader.ExecutionLog, error) {
	if logPath == "" {
		// Try to find logs alongside specs: spec.yaml -> spec.log.json
		var logs []*loader.ExecutionLog
		for _, sp := range specPaths {
			ext := filepath.Ext(sp)
			logFile := strings.TrimSuffix(sp, ext) + ".log.json"
			log, err := loader.LoadFile(logFile)
			if err != nil {
				return nil, fmt.Errorf("no --logs specified and cannot load %s: %w", logFile, err)
			}
			logs = append(logs, log)
		}
		return logs, nil
	}

	info, err := os.Stat(logPath)
	if err != nil {
		return nil, fmt.Errorf("cannot access logs at %s: %w", logPath, err)
	}

	if !info.IsDir() {
		// Single log file, apply to all specs
		log, err := loader.LoadFile(logPath)
		if err != nil {
			return nil, err
		}
		logs := make([]*loader.ExecutionLog, len(specPaths))
		for i := range logs {
			logs[i] = log
		}
		return logs, nil
	}

	// Directory of logs -- match by name
	var logs []*loader.ExecutionLog
	for _, sp := range specPaths {
		base := filepath.Base(sp)
		ext := filepath.Ext(base)
		logFile := filepath.Join(logPath, strings.TrimSuffix(base, ext)+".log.json")
		log, err := loader.LoadFile(logFile)
		if err != nil {
			return nil, fmt.Errorf("load log for %s: %w", sp, err)
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func validateSpecs(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot access %s: %w", path, err)
	}

	var specs []*spec.Spec
	var paths []string

	if info.IsDir() {
		specs, paths, err = spec.ParseDir(path)
		if err != nil {
			return err
		}
	} else {
		s, err := spec.ParseFile(path)
		if err != nil {
			return err
		}
		specs = []*spec.Spec{s}
		paths = []string{path}
	}

	if len(specs) == 0 {
		fmt.Println("No specs found.")
		return nil
	}

	valid := 0
	for i, s := range specs {
		if err := spec.Validate(s); err != nil {
			fmt.Fprintf(os.Stderr, "INVALID %s: %v\n", paths[i], err)
		} else {
			fmt.Printf("VALID   %s\n", paths[i])
			valid++
		}
	}
	fmt.Printf("\n%d/%d specs valid\n", valid, len(specs))

	if valid < len(specs) {
		os.Exit(1)
	}
	return nil
}

func initSpecs() error {
	dir := "specs"
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create specs dir: %w", err)
	}

	specContent := `name: "File reader agent reads requested files"
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
}
`

	specFile := filepath.Join(dir, "file-reader.yaml")
	logFile := filepath.Join(dir, "file-reader.log.json")

	if err := os.WriteFile(specFile, []byte(specContent), 0o644); err != nil {
		return fmt.Errorf("write spec: %w", err)
	}
	if err := os.WriteFile(logFile, []byte(logContent), 0o644); err != nil {
		return fmt.Errorf("write log: %w", err)
	}

	fmt.Printf("Created %s\n", specFile)
	fmt.Printf("Created %s\n", logFile)
	fmt.Println("\nRun: agentspec run specs/")
	return nil
}
