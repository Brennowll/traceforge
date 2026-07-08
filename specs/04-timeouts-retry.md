# Fase 4 — Timeouts e retry

## Objetivo

Adicionar comportamento realista de timeout e retry.

## Escopo

Cada chamada pode ter:

- `timeout_ms`
- `retry.attempts`
- `retry.backoff_ms`

## Regras

- Se `latency_ms > timeout_ms`, a tentativa vira timeout.
- Se a tentativa falhar, pode tentar novamente.
- Se a tentativa der timeout, pode tentar novamente.
- Backoff soma na latência total.
- Depois de esgotar retries, chamada falha.

## Exemplo

```yaml
services:
  api:
    calls:
      - service: payments
        timeout_ms: 200
        retry:
          attempts: 2
          backoff_ms: 50

  payments:
    latency_ms:
      min: 300
      max: 300
```

Resultado esperado:

```text
api -> payments
  attempt 1: timeout after 200ms
  backoff: 50ms
  attempt 2: timeout after 200ms
  backoff: 50ms
  attempt 3: timeout after 200ms

Result: failed
```

## Testes antes da implementação

Criar testes para:

- Timeout quando latência excede timeout.
- Sucesso quando latência é menor que timeout.
- Retry após falha.
- Retry após timeout.
- Backoff entra na latência total.
- Falha final após retries esgotados.
- Número correto de tentativas.

## Critérios de aceite

- Timeout funciona.
- Retry funciona.
- Backoff funciona.
- Trace explica cada tentativa.
- Testes determinísticos passam.

## Commit sugerido

```text
feat: add timeout and retry behavior
```
