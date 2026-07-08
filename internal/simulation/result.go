package simulation

import "github.com/brenno/traceforge/internal/stats"

// SimulationResult is the result of a single request simulation.
type SimulationResult struct {
	RequestID    int
	Success      bool
	TotalLatency int
	Trace        []TraceEvent
	Retries      int
	Timeouts     int
	Failures     int
	Error        string
}

// TraceEvent describes one event in the textual trace.
type TraceEvent struct {
	From      string
	To        string
	Attempt   int
	Status    string
	LatencyMS int
	Message   string
}

// BatchResult is the result of multiple request simulations.
type BatchResult struct {
	Results []SimulationResult
	Stats   stats.Stats
}
