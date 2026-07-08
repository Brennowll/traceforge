# Architecture

TraceForge is intentionally small and CLI-first.

## Packages

## `cmd/traceforge`

Cobra-based CLI with two commands:

- `traceforge version`
- `traceforge run scenario.yml --entry api`

The CLI is thin: it loads a scenario, validates the selected entry service, runs the simulator with the command context, renders text, and optionally writes HTML. The root command installs SIGINT/SIGTERM handling so long-running batches can be cancelled cleanly.

## `internal/scenario`

Owns YAML-facing data structures, parsing, validation, and defaults.

Responsibilities:

- Parse YAML using `gopkg.in/yaml.v3`.
- Validate service references, latency bounds, failure rates, timeouts, retries, cycles, and duplicate service names.
- Apply `simulation.default_timeout_ms` when a call omits `timeout_ms`.

## `internal/simulation`

Owns request execution.

Responsibilities:

- Execute calls depth-first from an entry service.
- Simulate latency and failures with an injected `RandomSource`.
- Apply timeout, retry, and backoff rules.
- Generate in-memory trace events.
- Run batches with a worker pool and fixed concurrency.
- Derive deterministic per-request seeds with `seed + requestNumber`.

## `internal/stats`

Owns aggregate batch statistics:

- total requests
- successes / failures
- success rate
- min / max / average latency
- p50 / p95 / p99
- total retries / timeouts

## `internal/output`

Owns rendering:

- terminal text output
- static HTML report with `html/template`

## Concurrency model

Batch simulation uses a worker pool. Each worker receives request indexes from a channel and writes its result into a preallocated slice at that index. This avoids append races and preserves deterministic output ordering. Workers check `context.Context` before dequeuing and during service traversal, returning cancellation errors to the CLI.

Each request gets its own random source. When `--seed` is set, request `N` uses seed `seed + N`, so `--concurrency 1` and `--concurrency 20` produce the same ordered results.

## Error handling

Scenario validation runs before simulation. Runtime errors are reserved for missing entrypoints, context cancellation, and filesystem/reporting failures.
