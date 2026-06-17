# Autentisering

## Oversikt

Skriveoperasjoner (`POST`, `PATCH`, `DELETE`) er beskyttet med JWT-autentisering via GitLab ID tokens.
Les-operasjoner (`GET`) og helsesjekk (`/healthz`) er åpne.

## GitLab ID tokens

GitLab 18.x utsteder automatisk signerte JWT-tokens til pipelines via ID tokens:

```yaml
# .gitlab-ci.yml
job:
  id_tokens:
    MY_TOKEN:
      aud: "https://din-app.com"
  script:
    - curl -X POST https://din-app/tasks
        -H "Authorization: Bearer $MY_TOKEN"
        -H "Content-Type: application/json"
        -d '{"title": "task fra pipeline"}'
```

## Verifiseringsflyt

```
1. Pipeline sender JWT i Authorization-header
2. Appen henter GitLab sin offentlige nøkkel fra JWKS-endepunkt
3. Appen verifiserer token-signaturen
4. Appen sjekker at aud-claim stemmer med forventet verdi
5. Appen sjekker at namespace_path er i AllowedNamespaces
6. Tillat eller avvis med 401/403
```

## JWT-claims fra GitLab

GitLab setter disse claims i ID tokens:

| Claim | Eksempel | Beskrivelse |
|-------|---------|-------------|
| `iss` | `https://gitlab.example.com` | Utsteder |
| `aud` | `https://din-app.com` | Mottaker |
| `namespace_path` | `min-gruppe` | GitLab-gruppe |
| `project_path` | `min-gruppe/mitt-repo` | Fullt prosjektpath |
| `ref` | `main` | Branch eller tag |
| `ref_type` | `branch` | Type referanse |

## Konfigurasjon

| Miljøvariabel | Beskrivelse | Eksempel |
|---------------|-------------|---------|
| `JWKS_URL` | GitLab sitt JWKS-endepunkt | `https://gitlab.example.com/-/jwks` |
| `JWT_AUDIENCE` | Forventet audience | `https://din-app.com` |
| `ALLOWED_NAMESPACES` | Tillatte GitLab-grupper (kommaseparert) | `gruppe-a,gruppe-b` |

`ALLOWED_NAMESPACES` bruker prefix-matching — `gruppe-a` tillater alle prosjekter under `gruppe-a/`.

## Feilresponser

| Statuskode | Årsak |
|------------|-------|
| 401 Unauthorized | Manglende token |
| 400 Bad Request | Ugyldig token eller manglende claims |
| 403 Forbidden | Gyldig token men namespace ikke tillatt |

## Lokal testing med mock

For lokal testing uten GitLab kan du bruke mock-verktøyene i `testdata/`:

```bash
# Generer RSA-nøkkelpar (kun én gang)
openssl genrsa -out testdata/private.pem 2048
openssl rsa -in testdata/private.pem -pubout -out testdata/public.pem

# Start mock JWKS-server
go run ./testdata/mock-jwks/

# Generer testtoken
TOKEN=$(go run ./testdata/gen-token/)

# Test beskyttet endepunkt
curl -s -X POST localhost:8072/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title": "test"}'
```

Start appen med mock-konfigurasjon:

```bash
JWKS_URL=http://localhost:8080/.well-known/jwks.json \
JWT_AUDIENCE=min-app-test \
ALLOWED_NAMESPACES=lagring \
STORE=memory \
go run main.go
```
