# Fase 8 — Relatório HTML estático

## Objetivo

Gerar um relatório HTML com resumo da simulação.

## Comando esperado

```bash
traceforge run scenario.yml --entry api --requests 100 --html report.html
```

## Conteúdo do relatório

- Nome do cenário.
- Entry service.
- Total de requests.
- Success rate.
- Latência média.
- p50, p95, p99.
- Total de retries.
- Total de timeouts.
- Amostra de traces.
- Tabela de requests.

## Regras

- HTML estático.
- Sem framework frontend.
- Sem assets externos obrigatórios.
- Template simples com `html/template`.

## Testes antes da implementação

Criar testes para:

- Gerar arquivo HTML.
- HTML contém estatísticas principais.
- HTML contém traces.
- Retornar erro se caminho inválido.

## Critérios de aceite

- `--html` gera arquivo.
- `report.html` abre no navegador.
- Sem dependência frontend.
- Testes passam.

## Commit sugerido

```text
feat: generate static html simulation report
```
