# Arkitektur

## Oversikt

Task Manager er en REST API bygget med Go sitt standardbibliotek. Applikasjonen følger en lagdelt arkitektur hvor hvert lag har ett klart ansvar.

```
HTTP-request
     │
     ▼
┌─────────────┐
│    Auth     │  JWT-validering for skriveoperasjoner
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Handler   │  Mottar HTTP, parser request, skriver response
└──────┬──────┘
       │ task.Task
       ▼
┌─────────────┐
│  TaskStore  │  Interface — definerer kontrakten
└──────┬──────┘
       │
   ┌───┴───┐
   │       │
   ▼       ▼
Postgres  Memory
Store     Store
```

## Pakkestruktur

### `internal/auth`

JWT-autentisering via GitLab ID tokens. Verifiserer token-signatur mot GitLab sitt JWKS-endepunkt og sjekker at `namespace_path` er i listen over tillatte grupper.

`jwt.go` — `Auth`-struct med `RequireJWT`-middleware:
```go
type Config struct {
    JWKSUrl           string
    Audience          string
    AllowedNamespaces []string
}
```

### `internal/task`

Kjernen i applikasjonen. Inneholder domenemodellen og kontrakten for lagring.

`task.go` — `Task`-structen med JSON-tags:
```go
type Task struct {
    ID        uint      `json:"id"`
    Title     string    `json:"title"`
    Done      bool      `json:"done"`
    CreatedAt time.Time `json:"created_at"`
}
```

`store.go` — `TaskStore`-interfacet som definerer hva en lagringsenhet må kunne gjøre, inkludert `Ping()` for helsesjekk.

`memory.go` — In-memory implementasjon med `map[uint]Task`. Brukes til utvikling og testing (`STORE=memory`).

### `internal/store`

PostgreSQL-implementasjon av `TaskStore`. Bruker `database/sql` med `pgx` som driver.
Kommuniserer med databasen via parameteriserte SQL-spørringer for å unngå SQL injection.

### `internal/handler`

HTTP-laget. `TaskHandler` tar imot en `TaskStore` via konstruktøren.
`RegisterRoutes` tar også imot en `*auth.Auth` for å beskytte skriveoperasjoner.

```
GET    /tasks        → handleGetAll      (åpen)
POST   /tasks        → handleSave        (krever JWT)
GET    /tasks/{id}   → handleGetById     (åpen)
PATCH  /tasks/{id}   → handleUpdate      (krever JWT)
DELETE /tasks/{id}   → handleDelete      (krever JWT)
GET    /healthz      → handleHealth      (åpen — brukes av k8s probes)
```

## Viktige designvalg

### Interface-drevet design

`TaskStore`-interfacet gjør at HTTP-laget er fullstendig uavhengig av lagringsimplementasjonen.
Å bytte fra in-memory til PostgreSQL krevde null endringer i `handler`-pakken — bare `main.go` endres.

### JWT som middleware

`RequireJWT` wrapper handler-funksjoner uten å endre dem. Dette holder autentiseringslogikken
adskilt fra forretningslogikken og gjør det enkelt å beskytte nye endepunkter.

### PATCH fremfor PUT

`handleUpdate` bruker PATCH og `UpdateTaskRequest` med peker-felt (`*string`, `*bool`).
Dette gir tre tilstander per felt: ikke sendt (`nil`), sendt tom, eller sendt med verdi.

### Konfigurasjon via miljøvariabler

All konfigurasjon leses fra miljøvariabler. Dette gjør applikasjonen portabel — samme binary
kjøres lokalt, i Docker og i Kubernetes uten kodeendringer.

## Dataflyt — POST /tasks (autentisert)

```
1. Klient sender POST /tasks med JWT i Authorization-header
2. RequireJWT henter og verifiserer JWT mot JWKS-endepunkt
3. namespace_path sjekkes mot AllowedNamespaces
4. handleSave leser og dekoder JSON til Task
5. CreatedAt settes av serveren
6. store.Save() lagrer i PostgreSQL med RETURNING
7. Klient mottar 201 Created med task-objektet
```
