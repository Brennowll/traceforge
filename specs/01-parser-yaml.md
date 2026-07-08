# Fase 1 — Parser de cenário YAML

## Objetivo

Ler um arquivo YAML e transformar em structs Go.

## Escopo

Implementar apenas parsing e validação básica. Ainda não executar simulação.

## Exemplo aceito

```yaml
services:
  api:
    calls:
      - service: payments
        timeout_ms: 200

  payments:
    latency_ms:
      min: 80
      max: 400
```

## Tipos esperados

```go
type Scenario struct {
    Services map[string]ServiceConfig `yaml:"services"`
}

type ServiceConfig struct {
    FailureRate float64      `yaml:"failure_rate"`
    Latency    LatencyConfig `yaml:"latency_ms"`
    Calls      []CallConfig  `yaml:"calls"`
}

type LatencyConfig struct {
    Min int `yaml:"min"`
    Max int `yaml:"max"`
}

type CallConfig struct {
    Service   string      `yaml:"service"`
    TimeoutMS int         `yaml:"timeout_ms"`
    Retry     RetryConfig `yaml:"retry"`
}

type RetryConfig struct {
    Attempts  int `yaml:"attempts"`
    BackoffMS int `yaml:"backoff_ms"`
}
```

## Testes antes da implementação

Criar testes para:

- Carregar YAML válido.
- Retornar erro para arquivo inexistente.
- Retornar erro para YAML inválido.
- Validar que existe pelo menos um serviço.
- Validar que chamadas apontam para serviços existentes.
- Validar `latency.min <= latency.max`.
- Validar `failure_rate` entre `0` e `1`.
- Validar `timeout_ms` não negativo.

## Critérios de aceite

- Parser lê `examples/basic.yml`.
- Validador detecta referências quebradas.
- Testes cobrem sucesso e erro.
- Nenhuma simulação implementada ainda.

## Commit sugerido

```text
feat: add scenario parser and validation
```
