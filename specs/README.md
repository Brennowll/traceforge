# TraceForge — Especificações do projeto

TraceForge é uma CLI em Go para simular sistemas distribuídos a partir de um arquivo YAML.

O usuário define serviços, chamadas entre serviços, latência, falhas, timeouts e retries. A ferramenta executa a simulação e gera trace textual, estatísticas e, depois, um relatório HTML.

> **TraceForge = Distributed system simulator written in Go**

## Como usar estas specs

- Executar **uma fase por vez**.
- Cada fase deve seguir TDD: escrever testes antes da implementação.
- Não implementar requisitos de fases futuras antes da fase atual estar verde.
- Atualizar o README do projeto a cada fase.

## Stack preferida

- Go
- Cobra
- `gopkg.in/yaml.v3`
- `context`
- `errgroup`
- `slog`
- `testing` nativo
- `testify` opcional
- Docker
- GitHub Actions

## Preferências técnicas

- Go + Cobra + yaml.v3 + slog + testing nativo.
- Evitar banco de dados no MVP.
- Foco em simulação, concorrência, CLI e testes.

## Princípios obrigatórios

- TDD obrigatório.
- Uma fase por vez.
- Não implementar fase futura antes da fase atual estar verde.
- Código pequeno.
- Sem abstração prematura.
- Sem banco de dados no MVP.
- Sem API HTTP no MVP.
- Sem interface web no MVP.
- CLI primeiro.
- Testes antes da funcionalidade.
- README atualizado a cada fase.

## Exemplo de entrada

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
    failure_rate: 0.2
    latency_ms:
      min: 80
      max: 400
    calls:
      - service: fraud_check
        timeout_ms: 150

  fraud_check:
    failure_rate: 0.1
    latency_ms:
      min: 40
      max: 120
```

## Exemplo de comando

```bash
traceforge run scenario.yml
```

## Exemplo de saída

```text
REQUEST 001

api -> payments
  attempt 1: timeout after 200ms
  attempt 2: success in 176ms

payments -> fraud_check
  attempt 1: success in 88ms

Result: success
Total latency: 421ms
Retries: 1
Timeouts: 1
Failures: 0
```

## Estrutura inicial esperada

```text
traceforge/
  cmd/
    traceforge/
      main.go

  internal/
    scenario/
      parser.go
      validator.go
      types.go

    simulation/
      simulator.go
      service.go
      trace.go
      result.go

    output/
      text.go
      html.go

    stats/
      stats.go

  examples/
    basic.yml
    retries.yml
    failures.yml

  .github/
    workflows/
      ci.yml

  Dockerfile
  docker-compose.yml
  Makefile
  README.md
  go.mod
  go.sum
```

## Ordem ideal de execução

1. [Fase 0 — Setup do projeto](./00-setup.md)
2. [Fase 1 — Parser de cenário YAML](./01-parser-yaml.md)
3. [Fase 2 — Simulação sequencial simples](./02-simulacao-sequencial.md)
4. [Fase 3 — Saída textual da CLI](./03-saida-textual.md)
5. [Fase 4 — Timeouts e retry](./04-timeouts-retry.md)
6. [Fase 5 — Múltiplas requests e estatísticas](./05-requests-estatisticas.md)
7. [Fase 6 — Concorrência controlada](./06-concorrencia.md)
8. [Fase 7 — Seed determinística](./07-seed-deterministica.md)
9. [Fase 8 — Relatório HTML estático](./08-relatorio-html.md)
10. [Fase 9 — Validações avançadas de cenário](./09-validacoes-avancadas.md)
11. [Fase 10 — Documentação e polish final](./10-documentacao-polish.md)
12. [Backlog opcional](./backlog.md)
