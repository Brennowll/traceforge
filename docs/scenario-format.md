# Scenario format

TraceForge scenarios are YAML files.

## Top-level fields

```yaml
simulation:
  max_depth: 10
  default_timeout_ms: 500

services:
  api: {}
```

## `simulation`

Optional global settings.

| Field | Description |
| --- | --- |
| `max_depth` | Allows cyclic service graphs up to this recursion depth. Without it, cycles are invalid. |
| `default_timeout_ms` | Timeout applied to calls that omit `timeout_ms`. |

## `services`

Required map of service name to service configuration.

```yaml
services:
  payments:
    failure_rate: 0.2
    latency_ms:
      min: 80
      max: 400
```

### Service fields

| Field | Description |
| --- | --- |
| `failure_rate` | Probability from `0` to `1` that a call to this service fails. Defaults to `0`. |
| `latency_ms.min` | Minimum simulated latency in milliseconds. Defaults to `0`. |
| `latency_ms.max` | Maximum simulated latency in milliseconds. Defaults to `0`. |
| `calls` | Outgoing calls performed by the service. |

## Calls

```yaml
calls:
  - service: payments
    timeout_ms: 200
    retry:
      attempts: 2
      backoff_ms: 50
```

| Field | Description |
| --- | --- |
| `service` | Target service name. Must exist in `services`. |
| `timeout_ms` | Required unless `simulation.default_timeout_ms` is set. `0` is valid and means an immediate timeout for any positive latency. |
| `retry.attempts` | Number of retries after the first attempt. `2` means up to 3 total attempts. |
| `retry.backoff_ms` | Latency added between failed/timed-out attempts. |

## Validation rules

- At least one service is required.
- Service names cannot be empty or duplicated.
- Calls must point to existing services.
- `latency_ms.min <= latency_ms.max`.
- Latency values must be non-negative.
- `failure_rate` must be between `0` and `1`.
- `timeout_ms` must be non-negative and present unless a default timeout is configured.
- `retry.attempts >= 0`.
- `retry.backoff_ms >= 0`.
- Cycles require `simulation.max_depth`.
- A valid graph may have only one unreferenced service, which is treated as the entry candidate.
- When an entry service is selected, all services must be reachable from that entry.

## Examples

See:

- `examples/basic.yml`
- `examples/failures.yml`
- `examples/retries.yml`
- `examples/concurrency.yml`
- `examples/invalid-commented.yml`
