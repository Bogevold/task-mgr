# Task Manager — Go Opplæringsprosjekt

Et fullstack REST API bygget i Go, med PostgreSQL backend, Docker, k3s deploy og JWT-autentisering.
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
│   ├── auth/           # JWT-autentisering
│   │   └── jwt.go
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
├── helm/               # Helm chart for k3s deploy
│   ├── Chart.yaml
│   ├── values.yaml
│   └── templates/
├── k8s/                # Kubernetes secrets (ikke i git)
├── testdata/           # Mock JWKS-server og token-generator
│   ├── mock-jwks/
│   ├── gen-token/
│   ├── public.pem
│   └── private.pem
├── docs/               # Dokumentasjon
├── Dockerfile
├── Dockerfile.migrate
├── docker-compose.yml
├── Makefile
├── go.mod
└── main.go
```

## Dokumentasjon

| Dokument | Beskrivelse |
|----------|-------------|
| [Arkitektur](docs/architecture.md) | Designvalg, pakkestruktur og dataflyt |
| [API](docs/api.md) | Endepunkter, request/response eksempler |
| [Autentisering](docs/auth.md) | JWT-autentisering med GitLab ID tokens |
| [Database](docs/database.md) | Skjema, migrasjoner og tilkoblingsoppsett |
| [Docker](docs/docker.md) | Dockerfile, docker-compose og lokalt oppsett |
| [Helm](docs/helm.md) | Helm chart struktur og deploy |
| [Kubernetes](docs/kubernetes.md) | k3s deploy, manifester og Traefik ingress |
| [Kommandoer](docs/commands.md) | Nyttige kommandoer for Go, Docker, kubectl, Helm og git |

## Teknologier

| Teknologi | Bruk |
|-----------|------|
| Go 1.25 | Applikasjonsspråk |
| `net/http` | HTTP-server (innebygd) |
| `database/sql` + pgx | PostgreSQL-driver |
| golang-migrate | Databasemigrasjoner |
| lestrrat-go/jwx/v2 | JWT-verifisering |
| PostgreSQL 16 | Databasebackend |
| Docker | Containerisering |
| Helm | Kubernetes pakkebehandler |
| k3s | Kubernetes-distribusjon |
| Traefik | Ingress-controller |

## Miljøvariabler

| Variabel | Beskrivelse | Standard |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL tilkoblingsstreng | — |
| `PORT` | Porten appen lytter på | `8072` |
| `STORE` | `postgres` eller `memory` (debug) | `postgres` |
| `JWKS_URL` | URL til JWKS-endepunkt | — |
| `JWT_AUDIENCE` | Forventet audience i JWT | — |
| `ALLOWED_NAMESPACES` | Kommaseparert liste over tillatte GitLab-grupper | — |
