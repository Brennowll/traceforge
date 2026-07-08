# Fase 2 — Simulação sequencial simples

## Objetivo

Executar uma request simulada de forma sequencial, sem concorrência real ainda.

## Escopo

Nesta fase, o simulador deve:

- Começar por um serviço raiz.
- Percorrer chamadas declaradas.
- Simular latência.
- Simular sucesso/falha.
- Gerar trace em memória.

## Comando esperado

```bash
traceforge run examples/basic.yml --entry api
```

## Modelo de resultado

```go
type SimulationResult struct {
    Success      bool
    TotalLatency int
    Trace        []TraceEvent
}

type TraceEvent struct {
    From       string
    To         string
    Attempt    int
    Status     string
    LatencyMS  int
    Message    string
}
```

## Testes antes da implementação

Criar testes para:

- Executar cenário sem chamadas.
- Executar cenário com uma chamada.
- Executar cenário com chamadas encadeadas.
- Registrar trace na ordem correta.
- Calcular latência total.
- Retornar erro se entrypoint não existir.

## Observação importante

Para testes determinísticos, não usar `math/rand` diretamente no domínio.

Criar uma interface:

```go
type RandomSource interface {
    Float64() float64
    IntBetween(min int, max int) int
}
```

Nos testes, usar um fake.

## Critérios de aceite

- Simulação sequencial funciona.
- Resultado é determinístico nos testes.
- Trace é gerado em memória.
- CLI executa cenário básico.

## Commit sugerido

```text
feat: add sequential simulation engine
```
