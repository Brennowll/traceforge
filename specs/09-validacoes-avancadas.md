# Fase 9 — Validações avançadas de cenário

## Objetivo

Evitar cenários inválidos ou perigosos.

## Validações

- Detectar serviços não referenciados.
- Detectar ciclos.
- Permitir ciclos apenas se `max_depth` existir.
- Validar timeout obrigatório quando há chamadas.
- Validar `retry.attempts >= 0`.
- Validar `retry.backoff_ms >= 0`.
- Validar nomes duplicados ou vazios.
- Validar latência não negativa.

## Configuração nova

```yaml
simulation:
  max_depth: 10
  default_timeout_ms: 500
```

## Testes antes da implementação

Criar testes para:

- Detectar ciclo sem `max_depth`.
- Permitir ciclo com `max_depth`.
- Aplicar `default_timeout_ms`.
- Detectar service name vazio.
- Detectar latency negativa.

## Critérios de aceite

- Cenários inválidos retornam mensagens claras.
- README documenta regras.
- Examples incluem cenário inválido comentado.

## Commit sugerido

```text
feat: add advanced scenario validation
```
