package stats

import (
	"math"
	"sort"
)

// ResultSummary is the minimal simulation data needed to calculate statistics.
type ResultSummary struct {
	Success      bool
	TotalLatency int
	Retries      int
	Timeouts     int
}

// Stats summarizes a batch of requests.
type Stats struct {
	TotalRequests int
	Successes     int
	Failures      int
	SuccessRate   float64
	MinLatency    int
	MaxLatency    int
	AvgLatency    float64
	P50           int
	P95           int
	P99           int
	TotalRetries  int
	TotalTimeouts int
}

// Calculate computes aggregate statistics for simulation results.
func Calculate(results []ResultSummary) Stats {
	out := Stats{TotalRequests: len(results)}
	if len(results) == 0 {
		return out
	}

	latencies := make([]int, len(results))
	sum := 0
	out.MinLatency = results[0].TotalLatency
	out.MaxLatency = results[0].TotalLatency
	for i, result := range results {
		if result.Success {
			out.Successes++
		}
		latencies[i] = result.TotalLatency
		sum += result.TotalLatency
		if result.TotalLatency < out.MinLatency {
			out.MinLatency = result.TotalLatency
		}
		if result.TotalLatency > out.MaxLatency {
			out.MaxLatency = result.TotalLatency
		}
		out.TotalRetries += result.Retries
		out.TotalTimeouts += result.Timeouts
	}
	out.Failures = out.TotalRequests - out.Successes
	out.SuccessRate = float64(out.Successes) / float64(out.TotalRequests) * 100
	out.AvgLatency = float64(sum) / float64(out.TotalRequests)
	sort.Ints(latencies)
	out.P50 = percentile(latencies, 50)
	out.P95 = percentile(latencies, 95)
	out.P99 = percentile(latencies, 99)
	return out
}

func percentile(sorted []int, p int) int {
	if len(sorted) == 0 {
		return 0
	}
	if p <= 0 {
		return sorted[0]
	}
	if p >= 100 {
		return sorted[len(sorted)-1]
	}
	idx := int(math.Ceil(float64(p)/100*float64(len(sorted)))) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}
