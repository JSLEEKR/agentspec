package runner

import (
	"fmt"
	"sync"

	"github.com/JSLEEKR/agentspec/internal/loader"
	"github.com/JSLEEKR/agentspec/internal/spec"
)

// RunParallel executes specs in parallel with the given concurrency limit.
func RunParallel(specs []*spec.Spec, paths []string, logs []*loader.ExecutionLog, workers int) (*RunResult, error) {
	if len(specs) != len(logs) {
		return nil, fmt.Errorf("spec count (%d) does not match log count (%d)", len(specs), len(logs))
	}
	if workers <= 0 {
		workers = 1
	}
	if workers > len(specs) {
		workers = len(specs)
	}

	type indexed struct {
		idx    int
		result SpecResult
	}

	results := make([]SpecResult, len(specs))
	ch := make(chan int, len(specs))
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range ch {
				sr := Run(specs[idx], logs[idx])
				if idx < len(paths) {
					sr.SpecPath = paths[idx]
				}
				results[idx] = sr
			}
		}()
	}

	for i := range specs {
		ch <- i
	}
	close(ch)
	wg.Wait()

	rr := &RunResult{Specs: results}
	for _, sr := range results {
		if sr.Passed {
			rr.TotalPassed++
		} else {
			rr.TotalFailed++
		}
	}

	return rr, nil
}
