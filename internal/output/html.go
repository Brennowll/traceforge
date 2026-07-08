package output

import (
	"html/template"
	"os"

	"github.com/brenno/traceforge/internal/simulation"
)

// HTMLReportData contains report metadata and simulation data.
type HTMLReportData struct {
	ScenarioName string
	Entry        string
	Batch        simulation.BatchResult
}

// HTMLRenderer writes static HTML reports.
type HTMLRenderer struct{}

// WriteFile writes a standalone HTML report.
func (HTMLRenderer) WriteFile(path string, data HTMLReportData) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return reportTemplate.Execute(file, data)
}

var reportTemplate = template.Must(template.New("report").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>TraceForge Report</title>
  <style>
    body { font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; margin: 2rem; color: #172033; }
    h1, h2 { color: #0f172a; }
    .cards { display: grid; grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)); gap: 1rem; margin: 1.5rem 0; }
    .card { border: 1px solid #d7dee8; border-radius: 10px; padding: 1rem; background: #f8fafc; }
    .value { font-size: 1.5rem; font-weight: 700; }
    table { border-collapse: collapse; width: 100%; margin: 1rem 0; }
    th, td { border: 1px solid #d7dee8; padding: .5rem; text-align: left; }
    th { background: #eef2f7; }
    pre { background: #0f172a; color: #e2e8f0; padding: 1rem; overflow: auto; border-radius: 10px; }
    .success { color: #047857; font-weight: 700; }
    .failed { color: #b91c1c; font-weight: 700; }
  </style>
</head>
<body>
  <h1>TraceForge Report</h1>
  <p><strong>Scenario:</strong> {{.ScenarioName}}</p>
  <p><strong>Entry service:</strong> {{.Entry}}</p>

  <h2>Summary</h2>
  <div class="cards">
    <div class="card"><div>Total requests</div><div class="value">{{.Batch.Stats.TotalRequests}}</div></div>
    <div class="card"><div>Success rate</div><div class="value">{{printf "%.2f" .Batch.Stats.SuccessRate}}%</div></div>
    <div class="card"><div>Avg latency</div><div class="value">{{printf "%.2f" .Batch.Stats.AvgLatency}}ms</div></div>
    <div class="card"><div>P50 / P95 / P99</div><div class="value">{{.Batch.Stats.P50}} / {{.Batch.Stats.P95}} / {{.Batch.Stats.P99}}ms</div></div>
    <div class="card"><div>Total retries</div><div class="value">{{.Batch.Stats.TotalRetries}}</div></div>
    <div class="card"><div>Total timeouts</div><div class="value">{{.Batch.Stats.TotalTimeouts}}</div></div>
  </div>

  <h2>Requests</h2>
  <table>
    <thead><tr><th>ID</th><th>Result</th><th>Latency</th><th>Retries</th><th>Timeouts</th><th>Failures</th></tr></thead>
    <tbody>
      {{range .Batch.Results}}
      <tr>
        <td>{{.RequestID}}</td>
        <td>{{if .Success}}<span class="success">success</span>{{else}}<span class="failed">failed</span>{{end}}</td>
        <td>{{.TotalLatency}}ms</td>
        <td>{{.Retries}}</td>
        <td>{{.Timeouts}}</td>
        <td>{{.Failures}}</td>
      </tr>
      {{end}}
    </tbody>
  </table>

  <h2>Trace samples</h2>
  {{range .Batch.Results}}
    {{if le .RequestID 5}}
    <h3>REQUEST {{printf "%03d" .RequestID}}</h3>
    <pre>{{range .Trace}}{{if ne .Status "backoff"}}{{.From}} -> {{.To}}
{{end}}{{if eq .Status "success"}}  attempt {{.Attempt}}: success in {{.LatencyMS}}ms
{{else if eq .Status "timeout"}}  attempt {{.Attempt}}: timeout after {{.LatencyMS}}ms
{{else if eq .Status "failure"}}  attempt {{.Attempt}}: failed in {{.LatencyMS}}ms
{{else if eq .Status "backoff"}}  backoff: {{.LatencyMS}}ms
{{end}}{{end}}</pre>
    {{end}}
  {{end}}
</body>
</html>`))
