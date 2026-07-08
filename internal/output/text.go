package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/brenno/traceforge/internal/simulation"
)

// TextRenderer renders simulation results for terminals.
type TextRenderer struct{}

// Render writes one or more request traces and statistics.
func (TextRenderer) Render(w io.Writer, batch simulation.BatchResult) error {
	for i, result := range batch.Results {
		if i > 0 {
			_, _ = fmt.Fprintln(w)
		}
		if err := renderResult(w, result); err != nil {
			return err
		}
	}
	if len(batch.Results) > 1 {
		_, _ = fmt.Fprintln(w)
		_, _ = fmt.Fprintln(w, "STATS")
		_, _ = fmt.Fprintf(w, "Total requests: %d\n", batch.Stats.TotalRequests)
		_, _ = fmt.Fprintf(w, "Successes: %d\n", batch.Stats.Successes)
		_, _ = fmt.Fprintf(w, "Failures: %d\n", batch.Stats.Failures)
		_, _ = fmt.Fprintf(w, "Success rate: %.2f%%\n", batch.Stats.SuccessRate)
		_, _ = fmt.Fprintf(w, "Min latency: %dms\n", batch.Stats.MinLatency)
		_, _ = fmt.Fprintf(w, "Max latency: %dms\n", batch.Stats.MaxLatency)
		_, _ = fmt.Fprintf(w, "Avg latency: %.2fms\n", batch.Stats.AvgLatency)
		_, _ = fmt.Fprintf(w, "P50: %dms\n", batch.Stats.P50)
		_, _ = fmt.Fprintf(w, "P95: %dms\n", batch.Stats.P95)
		_, _ = fmt.Fprintf(w, "P99: %dms\n", batch.Stats.P99)
		_, _ = fmt.Fprintf(w, "Total retries: %d\n", batch.Stats.TotalRetries)
		_, _ = fmt.Fprintf(w, "Total timeouts: %d\n", batch.Stats.TotalTimeouts)
	}
	return nil
}

// RenderString returns the text rendering as a string.
func (r TextRenderer) RenderString(batch simulation.BatchResult) string {
	var b strings.Builder
	_ = r.Render(&b, batch)
	return b.String()
}

func renderResult(w io.Writer, result simulation.SimulationResult) error {
	requestID := result.RequestID
	if requestID == 0 {
		requestID = 1
	}
	_, _ = fmt.Fprintf(w, "REQUEST %03d\n\n", requestID)

	currentFrom, currentTo := "", ""
	for _, event := range result.Trace {
		if event.Status != simulation.StatusBackoff && (event.From != currentFrom || event.To != currentTo) {
			if currentFrom != "" || currentTo != "" {
				_, _ = fmt.Fprintln(w)
			}
			currentFrom, currentTo = event.From, event.To
			_, _ = fmt.Fprintf(w, "%s -> %s\n", event.From, event.To)
		}
		switch event.Status {
		case simulation.StatusSuccess:
			_, _ = fmt.Fprintf(w, "  attempt %d: success in %dms\n", event.Attempt, event.LatencyMS)
		case simulation.StatusTimeout:
			_, _ = fmt.Fprintf(w, "  attempt %d: timeout after %dms\n", event.Attempt, event.LatencyMS)
		case simulation.StatusFailure:
			_, _ = fmt.Fprintf(w, "  attempt %d: failed in %dms\n", event.Attempt, event.LatencyMS)
		case simulation.StatusBackoff:
			_, _ = fmt.Fprintf(w, "  backoff: %dms\n", event.LatencyMS)
		default:
			_, _ = fmt.Fprintf(w, "  %s\n", event.Message)
		}
	}
	if len(result.Trace) > 0 {
		_, _ = fmt.Fprintln(w)
	}
	status := "failed"
	if result.Success {
		status = "success"
	}
	_, _ = fmt.Fprintf(w, "Result: %s\n", status)
	_, _ = fmt.Fprintf(w, "Total latency: %dms\n", result.TotalLatency)
	_, _ = fmt.Fprintf(w, "Retries: %d\n", result.Retries)
	_, _ = fmt.Fprintf(w, "Timeouts: %d\n", result.Timeouts)
	_, _ = fmt.Fprintf(w, "Failures: %d\n", result.Failures)
	return nil
}
