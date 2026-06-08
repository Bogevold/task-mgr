# API Dokumentasjon

Base URL lokalt: `http://localhost:8072`
Base URL k3s: `http://<domene>/task-mgr-api`

## Endepunkter

### Hent alle tasks

```
GET /tasks
```

**Response 200 OK**
```json
[
  {
    "id": 1,
    "title": "Min første task",
    "done": false,
    "created_at": "2026-06-08T20:33:28.360097Z"
  }
]
```

Tom liste returneres som `[]`, ikke `null`.

---

### Opprett ny task

```
POST /tasks
Content-Type: application/json
```

**Request body**
```json
{
  "title": "Ny task"
}
```

`created_at` settes av serveren og skal ikke sendes av klienten.

**Response 201 Created**
```json
{
  "id": 2,
  "title": "Ny task",
  "done": false,
  "created_at": "2026-06-08T20:35:00.123456Z"
}
```

---

### Hent én task

```
GET /tasks/{id}
```

**Response 200 OK**
```json
{
  "id": 1,
  "title": "Min første task",
  "done": false,
  "created_at": "2026-06-08T20:33:28.360097Z"
}
```

**Response 500** hvis ID ikke finnes.

---

### Oppdater task (delvis)

```
PATCH /tasks/{id}
Content-Type: application/json
```

Sender bare feltene som skal endres. Utelatte felt beholdes uendret.

**Request body — endre bare tittel**
```json
{
  "title": "Oppdatert tittel"
}
```

**Request body — marker som ferdig**
```json
{
  "done": true
}
```

**Request body — endre begge**
```json
{
  "title": "Ny tittel",
  "done": true
}
```

**Response 200 OK**
```json
{
  "id": 1,
  "title": "Oppdatert tittel",
  "done": false,
  "created_at": "2026-06-08T20:33:28.360097Z"
}
```

---

### Slett task

```
DELETE /tasks/{id}
```

**Response 204 No Content** — ingen body ved suksess.

---

## Feilresponser

Alle feil returneres som plain text med passende HTTP-statuskode.

| Statuskode | Betydning |
|------------|-----------|
| 400 Bad Request | Ugyldig JSON i request body |
| 204 No Content | Slett suksess |
| 500 Internal Server Error | Serverfeil eller ID ikke funnet |

---

## curl-eksempler

```bash
# Hent alle
curl -s localhost:8072/tasks | jq .

# Opprett
curl -s -X POST localhost:8072/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "min task"}' | jq .

# Hent én
curl -s localhost:8072/tasks/1 | jq .

# Oppdater
curl -s -X PATCH localhost:8072/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"done": true}' | jq .

# Slett
curl -s -X DELETE localhost:8072/tasks/1
```
