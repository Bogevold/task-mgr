# Database

## Skjema

```sql
CREATE TABLE tasks (
  id         SERIAL PRIMARY KEY,
  title      TEXT NOT NULL,
  done       BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

`SERIAL` — automatisk inkrementerende heltall, tilsvarer `uint` i Go.
`TIMESTAMPTZ` — tidspunkt med tidssone. Go sin `time.Time` mapper direkte til denne typen.

## Migrasjoner

Migrasjoner håndteres med [golang-migrate](https://github.com/golang-migrate/migrate).
Hver migrasjon består av to filer — én for å kjøre fremover (`up`) og én for å rulle tilbake (`down`).

```
migrations/
  000001_create_tasks.up.sql
  000001_create_tasks.down.sql
```

`golang-migrate` holder styr på hvilke migrasjoner som er kjørt i tabellen `schema_migrations`.

### Kjøre migrasjoner manuelt

```bash
# Fremover
migrate -path migrations \
  -database "postgres://taskuser:taskpass@localhost:5432/taskdb?sslmode=disable" up

# Rulle tilbake én
migrate -path migrations \
  -database "postgres://taskuser:taskpass@localhost:5432/taskdb?sslmode=disable" down 1

# Sjekk status
migrate -path migrations \
  -database "postgres://taskuser:taskpass@localhost:5432/taskdb?sslmode=disable" version
```

### Migrasjoner i k3s

I Kubernetes kjøres migrasjoner automatisk av en init-container (`kjor-migrasjoner`) ved hver deploy.
Init-containeren bruker imaget `task-mgr-migrate` og kjører før app-containeren starter.

## Tilkobling

Applikasjonen bruker `database/sql` med `pgx/v5` som driver.

Tilkoblingsstrengen settes via miljøvariabelen `DATABASE_URL`:
```
postgres://bruker:passord@host:port/database?sslmode=disable
```

| Miljø | Host |
|-------|------|
| Lokalt | `localhost` |
| docker-compose | `db` (tjenestenavn) |
| k3s | `postgres` (Kubernetes service-navn) |

## Koble til databasen direkte

**Lokalt / docker-compose:**
```bash
docker exec -it task-mgr-db-1 psql -U taskuser -d taskdb
```

**k3s:**
```bash
kubectl exec -it -n task-mgr postgres-0 -- psql -U taskuser -d taskdb
```

**Nyttige psql-kommandoer:**
```sql
\dt                  -- list tabeller
\d tasks             -- beskriv tasks-tabellen
SELECT * FROM tasks; -- hent alle tasks
```
