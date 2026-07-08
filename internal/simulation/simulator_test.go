package simulation

import (
	"context"
	"errors"
	"testing"

	"github.com/brenno/traceforge/internal/scenario"
)

type fakeRandom struct {
	floats []float64
	ints   []int
	fi     int
	ii     int
}

func (f *fakeRandom) Float64() float64 {
	if f.fi >= len(f.floats) {
		return 1
	}
	v := f.floats[f.fi]
	f.fi++
	return v
}

func (f *fakeRandom) IntBetween(min int, max int) int {
	if f.ii >= len(f.ints) {
		return min
	}
	v := f.ints[f.ii]
	f.ii++
	return v
}

func TestRunScenarioWithoutCalls(t *testing.T) {
	result, err := New(&scenario.Scenario{Services: map[string]scenario.ServiceConfig{"api": {}}}).Run(context.Background(), "api", &fakeRandom{})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Success || result.TotalLatency != 0 || len(result.Trace) != 0 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestRunScenarioWithOneCall(t *testing.T) {
	s := scenarioWithCalls([]scenario.CallConfig{{Service: "payments", TimeoutMS: 200}})
	result, err := New(s).Run(context.Background(), "api", &fakeRandom{ints: []int{120}, floats: []float64{0.9}})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Success || result.TotalLatency != 120 {
		t.Fatalf("unexpected result: %+v", result)
	}
	if got := result.Trace[0]; got.From != "api" || got.To != "payments" || got.Status != StatusSuccess || got.LatencyMS != 120 {
		t.Fatalf("unexpected trace: %+v", got)
	}
}

func TestRunScenarioWithChainedCallsTraceOrder(t *testing.T) {
	s := &scenario.Scenario{Services: map[string]scenario.ServiceConfig{
		"api":      {Calls: []scenario.CallConfig{{Service: "payments", TimeoutMS: 200}}},
		"payments": {Latency: scenario.LatencyConfig{Min: 80, Max: 80}, Calls: []scenario.CallConfig{{Service: "fraud", TimeoutMS: 200}}},
		"fraud":    {Latency: scenario.LatencyConfig{Min: 40, Max: 40}},
	}}
	result, err := New(s).Run(context.Background(), "api", &fakeRandom{ints: []int{80, 40}, floats: []float64{1, 1}})
	if err != nil {
		t.Fatal(err)
	}
	if result.TotalLatency != 120 || len(result.Trace) != 2 {
		t.Fatalf("unexpected result: %+v", result)
	}
	if result.Trace[0].From != "api" || result.Trace[1].From != "payments" {
		t.Fatalf("unexpected trace order: %+v", result.Trace)
	}
}

func TestRunEntrypointMissing(t *testing.T) {
	_, err := New(&scenario.Scenario{Services: map[string]scenario.ServiceConfig{"api": {}}}).Run(context.Background(), "missing", &fakeRandom{})
	if err == nil {
		t.Fatal("expected missing entry error")
	}
}

func TestTimeoutWhenLatencyExceedsTimeout(t *testing.T) {
	s := scenarioWithCalls([]scenario.CallConfig{{Service: "payments", TimeoutMS: 100}})
	result, err := New(s).Run(context.Background(), "api", &fakeRandom{ints: []int{150}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Success || result.Timeouts != 1 || result.TotalLatency != 100 || result.Trace[0].Status != StatusTimeout {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestExplicitZeroTimeoutIsEnforced(t *testing.T) {
	s := scenarioWithCalls([]scenario.CallConfig{{Service: "payments", TimeoutMS: 0, TimeoutSet: true}})
	result, err := New(s).Run(context.Background(), "api", &fakeRandom{ints: []int{1}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Success || result.Timeouts != 1 || result.TotalLatency != 0 || result.Trace[0].Status != StatusTimeout {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestSuccessWhenLatencyBelowTimeout(t *testing.T) {
	s := scenarioWithCalls([]scenario.CallConfig{{Service: "payments", TimeoutMS: 100}})
	result, err := New(s).Run(context.Background(), "api", &fakeRandom{ints: []int{99}, floats: []float64{1}})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Success || result.Timeouts != 0 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestRetryAfterFailure(t *testing.T) {
	s := scenarioWithCalls([]scenario.CallConfig{{Service: "payments", TimeoutMS: 100, Retry: scenario.RetryConfig{Attempts: 1, BackoffMS: 10}}})
	s.Services["payments"] = scenario.ServiceConfig{FailureRate: 0.5, Latency: scenario.LatencyConfig{Min: 20, Max: 20}}
	result, err := New(s).Run(context.Background(), "api", &fakeRandom{ints: []int{20, 20}, floats: []float64{0.1, 0.9}})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Success || result.Retries != 1 || result.TotalLatency != 50 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestRetryAfterTimeoutAndBackoff(t *testing.T) {
	s := scenarioWithCalls([]scenario.CallConfig{{Service: "payments", TimeoutMS: 100, Retry: scenario.RetryConfig{Attempts: 1, BackoffMS: 10}}})
	result, err := New(s).Run(context.Background(), "api", &fakeRandom{ints: []int{150, 50}, floats: []float64{1}})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Success || result.Retries != 1 || result.Timeouts != 1 || result.TotalLatency != 160 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestFinalFailureAfterRetries(t *testing.T) {
	s := scenarioWithCalls([]scenario.CallConfig{{Service: "payments", TimeoutMS: 100, Retry: scenario.RetryConfig{Attempts: 2, BackoffMS: 5}}})
	result, err := New(s).Run(context.Background(), "api", &fakeRandom{ints: []int{150, 150, 150}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Success || result.Retries != 2 || result.Timeouts != 3 || result.TotalLatency != 310 || len(result.Trace) != 5 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestRunBatchRunsNRequests(t *testing.T) {
	s := scenarioWithCalls([]scenario.CallConfig{{Service: "payments", TimeoutMS: 100}})
	batch, err := New(s).RunBatch(context.Background(), "api", BatchOptions{Requests: 5, Concurrency: 2, RandomFactory: func(int) RandomSource {
		return &fakeRandom{ints: []int{20}, floats: []float64{1}}
	}})
	if err != nil {
		t.Fatal(err)
	}
	if len(batch.Results) != 5 || batch.Stats.TotalRequests != 5 || batch.Stats.Successes != 5 {
		t.Fatalf("unexpected batch: %+v", batch)
	}
}

func TestRunBatchContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := New(&scenario.Scenario{Services: map[string]scenario.ServiceConfig{"api": {}}}).RunBatch(ctx, "api", BatchOptions{Requests: 1})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
}

func TestSeededBatchIsReproducible(t *testing.T) {
	s := scenarioWithCalls([]scenario.CallConfig{{Service: "payments", TimeoutMS: 100, Retry: scenario.RetryConfig{Attempts: 1, BackoffMS: 2}}})
	s.Services["payments"] = scenario.ServiceConfig{FailureRate: 0.3, Latency: scenario.LatencyConfig{Min: 1, Max: 90}}
	a, err := New(s).RunBatch(context.Background(), "api", BatchOptions{Requests: 50, Concurrency: 1, Seed: 42, Seeded: true})
	if err != nil {
		t.Fatal(err)
	}
	b, err := New(s).RunBatch(context.Background(), "api", BatchOptions{Requests: 50, Concurrency: 10, Seed: 42, Seeded: true})
	if err != nil {
		t.Fatal(err)
	}
	if !SameResults(a, b) {
		t.Fatal("same seed with different concurrency should produce same ordered results")
	}
	c, err := New(s).RunBatch(context.Background(), "api", BatchOptions{Requests: 50, Concurrency: 10, Seed: 99, Seeded: true})
	if err != nil {
		t.Fatal(err)
	}
	if SameResults(a, c) {
		t.Fatal("different seeds should alter results")
	}
}

func scenarioWithCalls(calls []scenario.CallConfig) *scenario.Scenario {
	return &scenario.Scenario{Services: map[string]scenario.ServiceConfig{
		"api":      {Calls: calls},
		"payments": {Latency: scenario.LatencyConfig{Min: 20, Max: 20}},
	}}
}
