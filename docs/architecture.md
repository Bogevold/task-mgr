# Arkitektur

## Oversikt

Task Manager er en REST API bygget med Go sitt standardbibliotek. Applikasjonen følger en lagdelt arkitektur hvor hvert lag har ett klart ansvar.

```
HTTP-request
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

`store.go` — `TaskStore`-interfacet som definerer hva en lagringsenhet må kunne gjøre.
Alle implementasjoner må oppfylle denne kontrakten.

`memory.go` — In-memory implementasjon med `map[uint]Task`. Brukes til utvikling og testing (`STORE=memory`).

### `internal/store`

PostgreSQL-implementasjon av `TaskStore`. Bruker `database/sql` med `pgx` som driver.
Kommuniserer med databasen via parameteriserte SQL-spørringer for å unngå SQL injection.

### `internal/handler`

HTTP-laget. `TaskHandler` tar imot en `TaskStore` via konstruktøren — den bryr seg ikke om det er Postgres eller in-memory bak.

`RegisterRoutes` kobler HTTP-metode og URL til riktig handler-funksjon:

```
GET    /tasks        → handleGetAll
POST   /tasks        → handleSave
GET    /tasks/{id}   → handleGetById
PATCH  /tasks/{id}   → handleUpdate
DELETE /tasks/{id}   → handleDelete
```

## Viktige designvalg

### Interface-drevet design

`TaskStore`-interfacet gjør at HTTP-laget er fullstendig uavhengig av lagringsimplementasjonen.
Å bytte fra in-memory til PostgreSQL krevde null endringer i `handler`-pakken — bare `main.go` endres.

### PATCH fremfor PUT

`handleUpdate` bruker PATCH og `UpdateTaskRequest` med peker-felt (`*string`, `*bool`).
Dette gir tre tilstander per felt: ikke sendt (`nil`), sendt tom, eller sendt med verdi.
PUT ville krevd at klienten sender alle felt og ville nullstilt utelatte felt.

### Feilhåndtering

Go returnerer feil som verdier, ikke exceptions. Alle feil håndteres eksplisitt med `if err != nil`.
HTTP-laget mapper feil til passende statuskoder: `400 Bad Request` for klientfeil, `500 Internal Server Error` for serverfeil.

### Konfigurasjon via miljøvariabler

All konfigurasjon (`DATABASE_URL`, `PORT`, `STORE`) leses fra miljøvariabler.
Dette gjør applikasjonen portabel — samme binary kjøres lokalt, i Docker og i Kubernetes
uten kodeendringer.

## Dataflyt — POST /tasks

```
1. Klient sender POST /tasks med JSON-body
2. handleSave leser og dekoder JSON til UpdateTaskRequest
3. CreatedAt settes av serveren (ikke klienten)
4. store.Save() lagrer i PostgreSQL med RETURNING
5. Lagret task (med tildelt ID) serialiseres til JSON
6. Klient mottar 201 Created med task-objektet
```
