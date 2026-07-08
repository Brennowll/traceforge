# Fase 7 — Seed determinística

## Objetivo

Permitir reproduzir simulações com uma seed.

## Comando esperado

```bash
traceforge run scenario.yml --entry api --requests 100 --seed 42
```

## Regras

- Mesma seed deve gerar mesmos resultados.
- Seeds diferentes podem gerar resultados diferentes.
- Em modo concorrente, evitar não determinismo excessivo.

## Implementação recomendada

Para simplificar:

- Cada request recebe uma seed derivada.
- Request 1 usa `seed + 1`.
- Request 2 usa `seed + 2`.
- Request N usa `seed + N`.

Assim a concorrência não muda o resultado estatístico.

## Testes antes da implementação

Criar testes para:

- Mesma seed gera mesmo batch.
- Seed diferente altera resultado.
- Concorrência não altera resultado final.

## Critérios de aceite

- `--seed` funciona.
- Simulações são reproduzíveis.
- Testes passam com e sem concorrência.

## Commit sugerido

```text
feat: add deterministic seeded simulations
```
