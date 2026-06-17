# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Task Manager REST API built in Go using only the standard library (`net/http` with Go 1.22+ pattern matching on `http.ServeMux`). PostgreSQL 16 for storage. Deployed to k3s via Helm charts.

## Common Commands

```bash
make test          # go test -v ./...
make vet           # go vet ./...
make up            # docker-compose up (app + postgres)
make down          # docker-compose down
make build         # build Docker app image
make build-migrate # build migration Docker image
make deploy        # helm upgrade to k3s
make ship          # build + push + deploy
make migrate       # run migrations locally via docker
```

Run a single test:
```bash
go test -v -run TestSaveAndGet ./internal/task/
```

## Architecture

Interface-driven layered design. The handler layer depends only on the `TaskStore` interface (`internal/task/store.go`), not concrete implementations.

- **`internal/task/`** — Domain model (`Task` struct), `TaskStore` interface, in-memory implementation
- **`internal/handler/`** — HTTP handlers and route registration. Uses pointer fields (`*string`, `*bool`) for PATCH partial updates
- **`internal/store/`** — PostgreSQL `TaskStore` implementation using `pgx/v5`
- **`internal/auth/`** — JWT/JWKS middleware. Protects write endpoints (POST/PATCH/DELETE); GET is open
- **`migrations/`** — SQL migrations for golang-migrate (run as k8s init-container)

Store selection is environment-based: `STORE=memory` for in-memory, `STORE=postgres` (default) for PostgreSQL.

## Environment Variables

| Variable | Purpose |
|---|---|
| `DATABASE_URL` | PostgreSQL connection string (required for postgres store) |
| `PORT` | Server listen port (default: 8072) |
| `STORE` | `postgres` or `memory` |
| `JWKS_URL` | JWKS endpoint for JWT validation |
| `JWT_AUDIENCE` | Expected JWT audience claim |
| `ALLOWED_NAMESPACES` | Comma-separated allowed namespace paths |

## Key Dependencies

- `github.com/jackc/pgx/v5` — PostgreSQL driver
- `github.com/lestrrat-go/jwx/v2` — JWT/JWKS parsing

No web framework — vanilla `net/http` throughout.

## Documentation

Detailed docs in `docs/` covering architecture, API reference, database, Docker, Kubernetes, and command reference.
