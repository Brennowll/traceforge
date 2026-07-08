# Fase 0 — Setup do projeto

## Objetivo

Criar o repositório Go com estrutura mínima, CI, Docker e comandos básicos.

## Tarefas

- Criar projeto Go.
- Criar estrutura de pastas.
- Criar Makefile.
- Criar Dockerfile.
- Criar `docker-compose.yml`.
- Criar GitHub Actions para rodar testes.
- Criar README inicial.
- Criar CLI mínima com comando `traceforge version`.

## Comandos esperados

```bash
make test
make lint
make build
make run
```

## Makefile mínimo

```makefile
.PHONY: test build run lint

test:
	go test ./...

build:
	go build -o bin/traceforge ./cmd/traceforge

run:
	go run ./cmd/traceforge

lint:
	go vet ./...
```

## Critérios de aceite

- `go test ./...` passa.
- `go vet ./...` passa.
- `docker build` funciona.
- GitHub Actions roda testes.
- `traceforge version` imprime a versão.

## Commit sugerido

```text
chore: setup go project structure
```
