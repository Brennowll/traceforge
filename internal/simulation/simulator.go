package simulation

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/brenno/traceforge/internal/scenario"
	"github.com/brenno/traceforge/internal/stats"
)

// Simulator executes scenarios.
type Simulator struct {
	Scenario *scenario.Scenario
}

// New creates a simulator for a parsed scenario.
func New(s *scenario.Scenario) *Simulator {
	return &Simulator{Scenario: s}
}

// BatchOptions controls batch execution.
type BatchOptions struct {
	Requests      int
	Concurrency   int
	Seed          int64
	Seeded        bool
	RandomFactory func(requestNumber int) RandomSource
}

// Run executes one request starting from entry.
func (s *Simulator) Run(ctx context.Context, entry string, rng RandomSource) (SimulationResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if s == nil || s.Scenario == nil {
		return SimulationResult{}, fmt.Errorf("simulation: scenario is nil")
	}
	if _, ok := s.Scenario.Services[entry]; !ok {
		return SimulationResult{}, fmt.Errorf("simulation: entry service %q does not exist", entry)
	}
	if rng == nil {
		rng = NewSeededRandom(time.Now().UnixNano())
	}
	result := SimulationResult{RequestID: 1, Success: true}
	ok, err := s.simulateService(ctx, entry, 0, rng, &result)
	if err != nil {
		return result, err
	}
	result.Success = ok
	if !ok && result.Failures == 0 {
		result.Failures = 1
	}
	return result, nil
}

func (s *Simulator) simulateService(ctx context.Context, name string, depth int, rng RandomSource, result *SimulationResult) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if maxDepth := s.Scenario.Simulation.MaxDepth; maxDepth > 0 && depth > maxDepth {
		result.Failures++
		result.Error = fmt.Sprintf("max depth %d reached at service %q", maxDepth, name)
		return false, nil
	}

	svc := s.Scenario.Services[name]
	for _, call := range svc.Calls {
		if err := ctx.Err(); err != nil {
			return false, err
		}
		callOK := s.executeCall(call, name, rng, result)
		if !callOK {
			result.Failures++
			return false, nil
		}
		childOK, err := s.simulateService(ctx, call.Service, depth+1, rng, result)
		if err != nil {
			return false, err
		}
		if !childOK {
			return false, nil
		}
	}
	return true, nil
}

func (s *Simulator) executeCall(call scenario.CallConfig, from string, rng RandomSource, result *SimulationResult) bool {
	target := s.Scenario.Services[call.Service]
	totalAttempts := call.Retry.Attempts + 1
	for attempt := 1; attempt <= totalAttempts; attempt++ {
		if attempt > 1 {
			result.Retries++
		}

		latency := rng.IntBetween(target.Latency.Min, target.Latency.Max)
		if callHasTimeout(call) && latency > call.TimeoutMS {
			result.Timeouts++
			result.TotalLatency += call.TimeoutMS
			result.addTrace(TraceEvent{
				From:      from,
				To:        call.Service,
				Attempt:   attempt,
				Status:    StatusTimeout,
				LatencyMS: call.TimeoutMS,
				Message:   fmt.Sprintf("timeout after %dms", call.TimeoutMS),
			})
			if attempt < totalAttempts {
				result.TotalLatency += call.Retry.BackoffMS
				result.addTrace(TraceEvent{From: from, To: call.Service, Status: StatusBackoff, LatencyMS: call.Retry.BackoffMS, Message: fmt.Sprintf("backoff: %dms", call.Retry.BackoffMS)})
				continue
			}
			return false
		}

		result.TotalLatency += latency
		if rng.Float64() < target.FailureRate {
			result.addTrace(TraceEvent{
				From:      from,
				To:        call.Service,
				Attempt:   attempt,
				Status:    StatusFailure,
				LatencyMS: latency,
				Message:   fmt.Sprintf("failed in %dms", latency),
			})
			if attempt < totalAttempts {
				result.TotalLatency += call.Retry.BackoffMS
				result.addTrace(TraceEvent{From: from, To: call.Service, Status: StatusBackoff, LatencyMS: call.Retry.BackoffMS, Message: fmt.Sprintf("backoff: %dms", call.Retry.BackoffMS)})
				continue
			}
			return false
		}

		result.addTrace(TraceEvent{
			From:      from,
			To:        call.Service,
			Attempt:   attempt,
			Status:    StatusSuccess,
			LatencyMS: latency,
			Message:   fmt.Sprintf("success in %dms", latency),
		})
		return true
	}
	return false
}

// RunBatch executes multiple requests, optionally in parallel.
func (s *Simulator) RunBatch(ctx context.Context, entry string, opts BatchOptions) (BatchResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if opts.Requests <= 0 {
		opts.Requests = 1
	}
	if opts.Concurrency <= 0 {
		opts.Concurrency = 1
	}
	if opts.Concurrency > opts.Requests {
		opts.Concurrency = opts.Requests
	}
	if s == nil || s.Scenario == nil {
		return BatchResult{}, fmt.Errorf("simulation: scenario is nil")
	}
	if _, ok := s.Scenario.Services[entry]; !ok {
		return BatchResult{}, fmt.Errorf("simulation: entry service %q does not exist", entry)
	}
	if err := ctx.Err(); err != nil {
		return BatchResult{}, err
	}

	results := make([]SimulationResult, opts.Requests)
	jobs := make(chan int)
	var wg sync.WaitGroup
	var firstErr error
	var errMu sync.Mutex
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	worker := func() {
		defer wg.Done()
		for idx := range jobs {
			if err := ctx.Err(); err != nil {
				errMu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				errMu.Unlock()
				return
			}
			rng := opts.randomSource(idx + 1)
			result, err := s.Run(ctx, entry, rng)
			result.RequestID = idx + 1
			if err != nil {
				errMu.Lock()
				if firstErr == nil {
					firstErr = err
					cancel()
				}
				errMu.Unlock()
				return
			}
			results[idx] = result
		}
	}

	wg.Add(opts.Concurrency)
	for i := 0; i < opts.Concurrency; i++ {
		go worker()
	}

	for i := 0; i < opts.Requests; i++ {
		select {
		case <-ctx.Done():
			errMu.Lock()
			if firstErr == nil {
				firstErr = ctx.Err()
			}
			errMu.Unlock()
			close(jobs)
			wg.Wait()
			return BatchResult{}, firstErr
		case jobs <- i:
		}
	}
	close(jobs)
	wg.Wait()

	errMu.Lock()
	err := firstErr
	errMu.Unlock()
	if err != nil {
		return BatchResult{}, err
	}

	return BatchResult{Results: results, Stats: stats.Calculate(toStatInputs(results))}, nil
}

func (o BatchOptions) randomSource(requestNumber int) RandomSource {
	if o.RandomFactory != nil {
		return o.RandomFactory(requestNumber)
	}
	if o.Seeded {
		return NewSeededRandom(o.Seed + int64(requestNumber))
	}
	return NewSeededRandom(time.Now().UnixNano() + int64(requestNumber))
}

func callHasTimeout(call scenario.CallConfig) bool {
	return call.TimeoutSet || call.TimeoutMS > 0
}

func toStatInputs(results []SimulationResult) []stats.ResultSummary {
	inputs := make([]stats.ResultSummary, len(results))
	for i, r := range results {
		inputs[i] = stats.ResultSummary{Success: r.Success, TotalLatency: r.TotalLatency, Retries: r.Retries, Timeouts: r.Timeouts}
	}
	return inputs
}

// SameResults reports whether two batches have identical ordered results.
func SameResults(a, b BatchResult) bool {
	if len(a.Results) != len(b.Results) {
		return false
	}
	for i := range a.Results {
		ar, br := a.Results[i], b.Results[i]
		if ar.Success != br.Success || ar.TotalLatency != br.TotalLatency || ar.Retries != br.Retries || ar.Timeouts != br.Timeouts || ar.Failures != br.Failures || len(ar.Trace) != len(br.Trace) {
			return false
		}
		for j := range ar.Trace {
			if ar.Trace[j] != br.Trace[j] {
				return false
			}
		}
	}
	return a.Stats == b.Stats
}

// SortResultsByRequestID sorts results by request id. It is mainly useful for callers that collect externally.
func SortResultsByRequestID(results []SimulationResult) {
	sort.Slice(results, func(i, j int) bool { return results[i].RequestID < results[j].RequestID })
}
