# Fase 6 — Concorrência controlada

## Objetivo

Executar múltiplas requests em paralelo, usando goroutines com limite de concorrência.

## Comando esperado

```bash
traceforge run scenario.yml --entry api --requests 1000 --concurrency 20
```

## Requisitos técnicos

- Usar `context.Context`.
- Usar worker pool ou `errgroup`.
- Respeitar limite de concorrência.
- Agregar resultados com segurança.
- Evitar data race.
- Testar com `go test -race`.

## Testes antes da implementação

Criar testes para:

- Rodar N requests concorrentes.
- Garantir quantidade correta de resultados.
- Respeitar cancelamento de contexto.
- Retornar erro em contexto cancelado.
- Não gerar data race.

## Critérios de aceite

- `--concurrency` funciona.
- `go test -race ./...` passa.
- Resultados são agregados corretamente.
- Context cancellation funciona.

## Commit sugerido

```text
feat: add concurrent batch simulation
```
