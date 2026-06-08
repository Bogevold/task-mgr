# Task Manager — Go Opplæringsprosjekt

Et fullstack REST API bygget i Go, med PostgreSQL backend, Docker og k3s deploy.
Prosjektet er brukt som læringsplattform for å utforske Go, REST, databaser, containers og Kubernetes.

## Kom i gang på 5 minutter

```bash
# Klon repoet
git clone https://github.com/bogevold/task-mgr
cd task-mgr

# Kopier og fyll inn secrets
cp k8s/secrets.example.yaml k8s/secret.yaml
# Rediger k8s/secret.yaml med ekte verdier

# Start lokalt med docker-compose
export DATABASE_URL="postgres://taskuser:taskpass@localhost:5432/taskdb?sslmode=disable"
docker compose up -d

# Kjør migrasjoner
migrate -path migrations -database "$DATABASE_URL" up

# Test
curl -s localhost:8072/tasks | jq .
```

## Prosjektstruktur

```
task-mgr/
├── internal/
│   ├── handler/        # HTTP-handlere og routing
│   │   └── task.go
│   ├── store/          # PostgreSQL implementasjon
│   │   └── postgres.go
│   └── task/           # Domenemodell, interface og in-memory store
│       ├── task.go
│       ├── store.go
│       ├── memory.go
│       └── memory_test.go
├── migrations/         # SQL migrasjoner (golang-migrate)
├── k8s/                # Kubernetes manifester
├── docs/               # Dokumentasjon
├── Dockerfile
├── Dockerfile.migrate
├── docker-compose.yml
├── go.mod
└── main.go
```

## Dokumentasjon

| Dokument | Beskrivelse |
|----------|-------------|
| [Arkitektur](docs/architecture.md) | Designvalg, pakkestruktur og dataflyt |
| [API](docs/api.md) | Endepunkter, request/response eksempler |
| [Database](docs/database.md) | Skjema, migrasjoner og tilkoblingsoppsett |
| [Docker](docs/docker.md) | Dockerfile, docker-compose og lokalt oppsett |
| [Kubernetes](docs/kubernetes.md) | k3s deploy, manifester og Traefik ingress |
| [Kommandoer](docs/commands.md) | Nyttige kommandoer for Go, Docker, kubectl og git |

## Teknologier

| Teknologi | Bruk |
|-----------|------|
| Go 1.25 | Applikasjonsspråk |
| `net/http` | HTTP-server (innebygd) |
| `database/sql` + pgx | PostgreSQL-driver |
| golang-migrate | Databasemigrasjoner |
| PostgreSQL 16 | Databasebackend |
| Docker | Containerisering |
| k3s | Kubernetes-distribusjon |
| Traefik | Ingress-controller |

## Miljøvariabler

| Variabel | Beskrivelse | Standard |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL tilkoblingsstreng | — |
| `PORT` | Porten appen lytter på | `8072` |
| `STORE` | `postgres` eller `memory` (debug) | `postgres` |
