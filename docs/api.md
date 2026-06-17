# API Dokumentasjon

Base URL lokalt: `http://localhost:8072`
Base URL k3s: `http://homelab/task-mgr-api`

## Autentisering

Skriveoperasjoner krever et gyldig GitLab ID token i `Authorization`-headeren:

```
Authorization: Bearer <token>
```

Se [Autentisering](auth.md) for detaljer.

## Endepunkter

### Hent alle tasks

```
GET /tasks
```

Åpent endepunkt — krever ikke autentisering.

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
Authorization: Bearer <token>
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

Åpent endepunkt — krever ikke autentisering.

**Response 200 OK**
```json
{
  "id": 1,
  "title": "Min første task",
  "done": false,
  "created_at": "2026-06-08T20:33:28.360097Z"
}
```

---

### Oppdater task (delvis)

```
PATCH /tasks/{id}
Authorization: Bearer <token>
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
Authorization: Bearer <token>
```

**Response 204 No Content** — ingen body ved suksess.

---

### Helsesjekk

```
GET /healthz
```

Åpent endepunkt — brukes av Kubernetes readiness og liveness probes.

**Response 200 OK**
```json
{"status": "ok"}
```

**Response 503 Service Unavailable** — databasen er ikke tilgjengelig.
```json
{"status": "error"}
```

---

## Feilresponser

| Statuskode | Betydning |
|------------|-----------|
| 400 Bad Request | Ugyldig JSON eller manglende/ugyldig token |
| 401 Unauthorized | Manglende token |
| 403 Forbidden | Gyldig token men ikke tillatt namespace |
| 204 No Content | Slett suksess |
| 500 Internal Server Error | Serverfeil eller ID ikke funnet |
| 503 Service Unavailable | Database utilgjengelig |
