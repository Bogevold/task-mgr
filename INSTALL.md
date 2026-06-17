# Installasjon av dokumentasjon

## Struktur

Dokumentasjonspakken inneholder:

```
README.md          → erstatter eksisterende README.md i roten
INSTALL.md         → denne filen (kan slettes etter installasjon)
docs/
  architecture.md  → ny eller erstatter eksisterende
  api.md           → ny eller erstatter eksisterende
  auth.md          → NY — ikke i forrige versjon
  database.md      → ny eller erstatter eksisterende
  docker.md        → ny eller erstatter eksisterende
  helm.md          → NY — ikke i forrige versjon
  kubernetes.md    → ny eller erstatter eksisterende
  commands.md      → ny eller erstatter eksisterende
```

## Installasjon

### Alternativ 1 — Pakk ut direkte i prosjektmappen (anbefalt)

```bash
# Stå i roten av task-mgr prosjektet
cd ~/Programmering/golang/task-mgr

# Pakk ut og overskriv eksisterende filer
unzip -o task-mgr-docs.zip

# Slett installasjonsfilen
rm INSTALL.md
```

`-o` flagget overskriv eksisterende filer uten å spørre.

### Alternativ 2 — Pakk ut et annet sted og kopier manuelt

```bash
# Pakk ut i en midlertidig mappe
mkdir /tmp/docs-update
unzip task-mgr-docs.zip -d /tmp/docs-update

# Kopier til prosjektet
cp /tmp/docs-update/README.md ~/Programmering/golang/task-mgr/
cp -r /tmp/docs-update/docs/ ~/Programmering/golang/task-mgr/

# Rydd opp
rm -rf /tmp/docs-update
```

## Etter installasjon

Commit dokumentasjonen:

```bash
gaa && gcam "docs: oppdatert dokumentasjon med auth, Helm og Makefile"
```

## Hva er nytt i denne versjonen

- `docs/auth.md` — JWT-autentisering med GitLab ID tokens og lokal mock-testing
- `docs/helm.md` — Helm chart struktur, verdier og deploy-kommandoer
- `docs/commands.md` — utvidet med Helm, Makefile og JWT-testing
- `docs/api.md` — oppdatert med autentiseringskrav og `/healthz`-endepunkt
- `docs/architecture.md` — oppdatert med auth-lag og ny endepunktstruktur
- `README.md` — oppdatert prosjektstruktur og miljøvariabler
