package stats

import "testing"

func TestCalculateStats(t *testing.T) {
	got := Calculate([]ResultSummary{
		{Success: true, TotalLatency: 10, Retries: 1, Timeouts: 0},
		{Success: false, TotalLatency: 20, Retries: 2, Timeouts: 1},
		{Success: true, TotalLatency: 30, Retries: 0, Timeouts: 2},
		{Success: true, TotalLatency: 40, Retries: 1, Timeouts: 0},
	})
	if got.TotalRequests != 4 || got.Successes != 3 || got.Failures != 1 {
		t.Fatalf("unexpected counts: %+v", got)
	}
	if got.SuccessRate != 75 || got.AvgLatency != 25 {
		t.Fatalf("unexpected rate/avg: %+v", got)
	}
	if got.MinLatency != 10 || got.MaxLatency != 40 {
		t.Fatalf("unexpected min/max: %+v", got)
	}
	if got.P50 != 20 || got.P95 != 40 || got.P99 != 40 {
		t.Fatalf("unexpected percentiles: %+v", got)
	}
	if got.TotalRetries != 4 || got.TotalTimeouts != 3 {
		t.Fatalf("unexpected totals: %+v", got)
	}
}

func TestCalculateEmptyStats(t *testing.T) {
	got := Calculate(nil)
	if got.TotalRequests != 0 || got.SuccessRate != 0 {
		t.Fatalf("unexpected empty stats: %+v", got)
	}
}
