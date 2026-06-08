# Docker

## Dockerfile

Applikasjonen bruker et multi-stage bygg for å holde det endelige imaget minimalt.

```
Stage 1 (builder): golang:1.25-alpine
  - Laster ned avhengigheter (go mod download)
  - Kompilerer Go-koden til en statisk binær

Stage 2 (runtime): alpine:3.22
  - Kopierer bare den kompilerte binæren fra stage 1
  - Resulterer i et lite image (~20MB vs ~300MB med Go-toolchain)
```

### Hvorfor alpine?

`alpine` er en minimal Linux-distribusjon (~5MB). Den inneholder nok til å kjøre
en Go-binær, men lite annet. Alternativet er `scratch` (helt tomt image) som er enda
mindre men mangler nyttige verktøy som shell og SSL-sertifikater.

## Dockerfile.migrate

Et eget image for å kjøre databasemigrasjoner. Inneholder bare:
- `migrate`-binæren
- SQL-migrasjonsfiler fra `migrations/`

Brukes som init-container i Kubernetes.

## docker-compose

Lokalt utviklingsmiljø med to tjenester:

```yaml
services:
  app:   # Go-applikasjonen
  db:    # PostgreSQL 16
```

`app` venter på at `db` er klar via `depends_on` med `condition: service_healthy`.
`db` har en healthcheck som bruker `pg_isready`.

### Starte og stoppe

```bash
# Start i bakgrunnen
docker compose up -d

# Se logger
docker compose logs -f

# Se logger for én tjeneste
docker compose logs -f app

# Stopp
docker compose down

# Stopp og slett volumes (mister databasedata)
docker compose down -v
```

## Bygge og tagge images

```bash
# Bygg app-image
docker build -t task-mgr .

# Bygg migrate-image
docker build -f Dockerfile.migrate -t task-mgr-migrate .

# Tag for lokalt registry
docker tag task-mgr homelab:30500/task-mgr:latest
docker tag task-mgr-migrate homelab:30500/task-mgr-migrate:latest

# Push til registry
docker push homelab:30500/task-mgr:latest
docker push homelab:30500/task-mgr-migrate:latest
```

## Miljøvariabler

| Variabel | Beskrivelse |
|----------|-------------|
| `DATABASE_URL` | PostgreSQL tilkoblingsstreng |
| `PORT` | Porten appen lytter på (standard: 8072) |
| `STORE` | `postgres` eller `memory` |

### Kjøre med in-memory store (debug)

```bash
docker run -e STORE=memory -e PORT=8072 -p 8072:8072 task-mgr
```
