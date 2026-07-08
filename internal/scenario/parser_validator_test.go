package scenario

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseValidYAML(t *testing.T) {
	s, err := Parse([]byte(`services:
  api:
    calls:
      - service: payments
        timeout_ms: 200
  payments:
    latency_ms:
      min: 80
      max: 400
`))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(s.Services) != 2 || s.Services["payments"].Latency.Min != 80 {
		t.Fatalf("unexpected scenario: %+v", s)
	}
}

func TestLoadFileMissing(t *testing.T) {
	_, err := LoadFile(filepath.Join(t.TempDir(), "missing.yml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestParseInvalidYAML(t *testing.T) {
	_, err := Parse([]byte("services: ["))
	if err == nil {
		t.Fatal("expected YAML error")
	}
}

func TestValidateRequiresService(t *testing.T) {
	_, err := Parse([]byte(`services: {}`))
	assertErrContains(t, err, "at least one service")
}

func TestValidateUnknownCallTarget(t *testing.T) {
	_, err := Parse([]byte(`services:
  api:
    calls:
      - service: missing
        timeout_ms: 1
`))
	assertErrContains(t, err, "unknown service")
}

func TestValidateLatencyMinMax(t *testing.T) {
	_, err := Parse([]byte(`services:
  api:
    latency_ms:
      min: 20
      max: 10
`))
	assertErrContains(t, err, "min must be <= max")
}

func TestValidateFailureRate(t *testing.T) {
	_, err := Parse([]byte(`services:
  api:
    failure_rate: 1.5
`))
	assertErrContains(t, err, "failure_rate")
}

func TestValidateTimeoutNonNegative(t *testing.T) {
	_, err := Parse([]byte(`services:
  api:
    calls:
      - service: payments
        timeout_ms: -1
  payments: {}
`))
	assertErrContains(t, err, "timeout_ms")
}

func TestValidateExplicitZeroTimeout(t *testing.T) {
	s, err := Parse([]byte(`services:
  api:
    calls:
      - service: payments
        timeout_ms: 0
  payments: {}
`))
	if err != nil {
		t.Fatalf("explicit timeout_ms: 0 should be valid: %v", err)
	}
	call := s.Services["api"].Calls[0]
	if call.TimeoutMS != 0 || !call.TimeoutSet {
		t.Fatalf("unexpected timeout state: %+v", call)
	}
}

func TestValidateRequiresTimeoutWhenOmitted(t *testing.T) {
	_, err := Parse([]byte(`services:
  api:
    calls:
      - service: payments
  payments: {}
`))
	assertErrContains(t, err, "requires timeout_ms")
}

func TestParseRejectsUnknownRetryField(t *testing.T) {
	_, err := Parse([]byte(`services:
  api:
    calls:
      - service: payments
        timeout_ms: 1
        retry:
          bogus: 1
  payments: {}
`))
	assertErrContains(t, err, "field bogus")
}

func TestLoadExamplesBasic(t *testing.T) {
	_, err := LoadFile(filepath.Join("..", "..", "examples", "basic.yml"))
	if err != nil {
		t.Fatalf("LoadFile examples/basic.yml: %v", err)
	}
}

func TestAdvancedValidationCycleRequiresMaxDepth(t *testing.T) {
	_, err := Parse([]byte(`services:
  api:
    calls:
      - service: payments
        timeout_ms: 100
  payments:
    calls:
      - service: api
        timeout_ms: 100
`))
	assertErrContains(t, err, "cycle detected")
}

func TestAdvancedValidationAllowsCycleWithMaxDepth(t *testing.T) {
	_, err := Parse([]byte(`simulation:
  max_depth: 3
services:
  api:
    calls:
      - service: payments
        timeout_ms: 100
  payments:
    calls:
      - service: api
        timeout_ms: 100
`))
	if err != nil {
		t.Fatalf("expected cycle with max_depth to be allowed: %v", err)
	}
}

func TestAdvancedValidationAppliesDefaultTimeout(t *testing.T) {
	s, err := Parse([]byte(`simulation:
  default_timeout_ms: 500
services:
  api:
    calls:
      - service: payments
  payments: {}
`))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if got := s.Services["api"].Calls[0].TimeoutMS; got != 500 {
		t.Fatalf("timeout = %d, want 500", got)
	}
}

func TestAdvancedValidationEmptyServiceName(t *testing.T) {
	_, err := Parse([]byte("services:\n  \"\": {}\n"))
	assertErrContains(t, err, "service name cannot be empty")
}

func TestAdvancedValidationNegativeLatency(t *testing.T) {
	_, err := Parse([]byte(`services:
  api:
    latency_ms:
      min: -1
      max: 1
`))
	assertErrContains(t, err, "latency must be non-negative")
}

func TestAdvancedValidationDuplicateServiceName(t *testing.T) {
	_, err := Parse([]byte(`services:
  api: {}
  api: {}
`))
	assertErrContains(t, err, "duplicate service name")
}

func TestValidateDetectsMultipleUnreferencedServices(t *testing.T) {
	_, err := Parse([]byte(`services:
  api: {}
  unused: {}
`))
	assertErrContains(t, err, "unreferenced")
}

func TestValidateEntryDetectsUnreachableCycle(t *testing.T) {
	s, err := Parse([]byte(`simulation:
  max_depth: 3
services:
  api:
    calls:
      - service: payments
        timeout_ms: 100
  payments: {}
  orphan_a:
    calls:
      - service: orphan_b
        timeout_ms: 100
  orphan_b:
    calls:
      - service: orphan_a
        timeout_ms: 100
`))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	err = ValidateEntry(s, "api")
	assertErrContains(t, err, "not reachable")
}

func TestLoadFileRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "scenario.yml")
	if err := os.WriteFile(path, []byte("services:\n  api: {}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadFile(path); err != nil {
		t.Fatalf("LoadFile returned error: %v", err)
	}
}

func assertErrContains(t *testing.T, err error, want string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error containing %q", want)
	}
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("error %q does not contain %q", err.Error(), want)
	}
}
