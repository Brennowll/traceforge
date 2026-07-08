# Fase 5 — Múltiplas requests e estatísticas

## Objetivo

Permitir rodar várias requests e calcular estatísticas.

## Comando esperado

```bash
traceforge run scenario.yml --entry api --requests 100
```

## Estatísticas

- Total requests.
- Successes.
- Failures.
- Success rate.
- Min latency.
- Max latency.
- Avg latency.
- P50.
- P95.
- P99.
- Total retries.
- Total timeouts.

## Tipos sugeridos

```go
type BatchResult struct {
    Results []SimulationResult
    Stats   Stats
}

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
```

## Testes antes da implementação

Criar testes para:

- Calcular success rate.
- Calcular média.
- Calcular p50.
- Calcular p95.
- Calcular p99.
- Somar retries.
- Somar timeouts.
- Rodar N requests.

## Critérios de aceite

- `--requests` funciona.
- Estatísticas são impressas.
- Testes de percentil passam.
- Código de stats isolado em `internal/stats`.

## Commit sugerido

```text
feat: add batch simulation statistics
```
