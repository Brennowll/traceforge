package output

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/brenno/traceforge/internal/simulation"
	"github.com/brenno/traceforge/internal/stats"
)

func TestTextRendererRendersSuccess(t *testing.T) {
	text := (TextRenderer{}).RenderString(batchWithResult(simulation.SimulationResult{RequestID: 1, Success: true, TotalLatency: 120, Trace: []simulation.TraceEvent{{From: "api", To: "payments", Attempt: 1, Status: simulation.StatusSuccess, LatencyMS: 120}}}))
	for _, want := range []string{"REQUEST 001", "api -> payments", "attempt 1: success in 120ms", "Result: success"} {
		if !strings.Contains(text, want) {
			t.Fatalf("text missing %q:\n%s", want, text)
		}
	}
}

func TestTextRendererRendersFailureAndMultipleEvents(t *testing.T) {
	text := (TextRenderer{}).RenderString(batchWithResult(simulation.SimulationResult{RequestID: 1, Success: false, TotalLatency: 250, Retries: 1, Timeouts: 1, Failures: 1, Trace: []simulation.TraceEvent{
		{From: "api", To: "payments", Attempt: 1, Status: simulation.StatusTimeout, LatencyMS: 100},
		{From: "api", To: "payments", Status: simulation.StatusBackoff, LatencyMS: 50},
		{From: "api", To: "payments", Attempt: 2, Status: simulation.StatusFailure, LatencyMS: 100},
	}}))
	for _, want := range []string{"attempt 1: timeout after 100ms", "backoff: 50ms", "attempt 2: failed in 100ms", "Result: failed", "Retries: 1", "Timeouts: 1", "Failures: 1"} {
		if !strings.Contains(text, want) {
			t.Fatalf("text missing %q:\n%s", want, text)
		}
	}
}

func TestTextRendererRendersBatchStats(t *testing.T) {
	batch := simulation.BatchResult{Results: []simulation.SimulationResult{{RequestID: 1, Success: true}, {RequestID: 2, Success: false}}, Stats: stats.Stats{TotalRequests: 2, Successes: 1, Failures: 1, SuccessRate: 50, AvgLatency: 10, P50: 10, P95: 10, P99: 10}}
	text := (TextRenderer{}).RenderString(batch)
	if !strings.Contains(text, "STATS") || !strings.Contains(text, "Success rate: 50.00%") {
		t.Fatalf("stats missing:\n%s", text)
	}
}

func TestHTMLRendererWritesReport(t *testing.T) {
	path := filepath.Join(t.TempDir(), "report.html")
	batch := batchWithResult(simulation.SimulationResult{RequestID: 1, Success: true, TotalLatency: 42, Retries: 1, Timeouts: 0, Trace: []simulation.TraceEvent{{From: "api", To: "payments", Attempt: 1, Status: simulation.StatusSuccess, LatencyMS: 42}}})
	batch.Stats = stats.Stats{TotalRequests: 1, Successes: 1, SuccessRate: 100, AvgLatency: 42, P50: 42, P95: 42, P99: 42, TotalRetries: 1}
	if err := (HTMLRenderer{}).WriteFile(path, HTMLReportData{ScenarioName: "basic.yml", Entry: "api", Batch: batch}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	html := string(data)
	for _, want := range []string{"TraceForge Report", "basic.yml", "Success rate", "100.00%", "api -> payments"} {
		if !strings.Contains(html, want) {
			t.Fatalf("html missing %q", want)
		}
	}
}

func TestHTMLRendererInvalidPath(t *testing.T) {
	err := (HTMLRenderer{}).WriteFile(filepath.Join(t.TempDir(), "missing", "report.html"), HTMLReportData{})
	if err == nil {
		t.Fatal("expected invalid path error")
	}
}

func batchWithResult(result simulation.SimulationResult) simulation.BatchResult {
	return simulation.BatchResult{Results: []simulation.SimulationResult{result}, Stats: stats.Stats{TotalRequests: 1, Successes: 1, SuccessRate: 100, AvgLatency: float64(result.TotalLatency), P50: result.TotalLatency, P95: result.TotalLatency, P99: result.TotalLatency, TotalRetries: result.Retries, TotalTimeouts: result.Timeouts}}
}
