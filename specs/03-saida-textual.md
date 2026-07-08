# Fase 3 — Saída textual da CLI

## Objetivo

Transformar o resultado da simulação em uma saída legível no terminal.

## Escopo

Implementar apenas saída texto.

## Exemplo

```text
REQUEST 001

api -> payments
  attempt 1: success in 120ms

Result: success
Total latency: 120ms
Retries: 0
Timeouts: 0
Failures: 0
```

## Tarefas

- Criar pacote `internal/output`.
- Implementar `TextRenderer`.
- Adicionar flags de CLI.
- Exibir resultado no terminal.

## Flags

```bash
traceforge run scenario.yml --entry api
```

## Testes antes da implementação

Criar testes para:

- Renderizar sucesso.
- Renderizar falha.
- Renderizar múltiplos eventos.
- Renderizar estatísticas básicas.

## Critérios de aceite

- CLI imprime trace legível.
- Testes de renderização passam.
- `examples/basic.yml` funciona via terminal.

## Commit sugerido

```text
feat: render simulation trace as text
```
