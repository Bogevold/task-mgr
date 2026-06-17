# Kommandoreferanse

## Go

```bash
# Bygg
go build ./...                    # Bygg alle pakker
go build -o task-mgr .            # Bygg og navngi binæren

# Kjør
go run main.go                    # Kjør direkte
STORE=memory go run main.go       # Kjør med in-memory store
PORT=9000 go run main.go          # Kjør på annen port

# Kjør med JWT-autentisering (lokal mock)
JWKS_URL=http://localhost:8080/.well-known/jwks.json \
JWT_AUDIENCE=min-app-test \
ALLOWED_NAMESPACES=lagring \
STORE=memory \
go run main.go

# Test
go test ./...                     # Kjør alle tester
go test -v ./...                  # Verbose output
go test ./internal/task/...       # Test én pakke

# Avhengigheter
go mod tidy                       # Rydd opp ubrukte avhengigheter
go get <pakke>                    # Legg til avhengighet
go mod download                   # Last ned alle avhengigheter

# Verktøy
go vet ./...                      # Statisk analyse
gofmt -w .                        # Formater kode
```

## Migrasjoner (golang-migrate)

```bash
# Kjør alle migrasjoner fremover
migrate -path migrations \
  -database "$DATABASE_URL" up

# Rulle tilbake én migrasjon
migrate -path migrations \
  -database "$DATABASE_URL" down 1

# Sjekk gjeldende versjon
migrate -path migrations \
  -database "$DATABASE_URL" version

# Lag ny migrasjon
migrate create -ext sql -dir migrations -seq <navn>
```

## Docker

```bash
# Bygg images
docker build -t task-mgr .
docker build -f Dockerfile.migrate -t task-mgr-migrate .

# Tag og push til lokalt registry
docker tag task-mgr homelab:30500/task-mgr:latest
docker tag task-mgr-migrate homelab:30500/task-mgr-migrate:latest
docker push homelab:30500/task-mgr:latest
docker push homelab:30500/task-mgr-migrate:latest

# docker-compose
docker compose up -d              # Start i bakgrunnen
docker compose down               # Stopp
docker compose down -v            # Stopp og slett volumes
docker compose logs -f            # Følg logger
docker compose logs -f app        # Følg logger for én tjeneste
docker compose ps                 # Se status
docker compose build              # Bygg images på nytt

# Koble til kjørende container
docker exec -it task-mgr-db-1 psql -U taskuser -d taskdb
docker exec -it task-mgr-app-1 sh
```

## Makefile

```bash
make help                         # Se alle tilgjengelige kommandoer
make version                      # Vis versjon, branch og image-tag

# Versjonering
make bump-patch                   # Bump patch (v0.0.x)
make bump-minor                   # Bump minor (v0.x.0)
make bump-major                   # Bump major (vx.0.0)

# Branching
make branch-fix NAME=beskrivelse      # Bump patch og lag fix/beskrivelse
make branch-feature NAME=beskrivelse  # Bump minor og lag feature/beskrivelse
make branch-major NAME=beskrivelse    # Bump major og lag major/beskrivelse

# Bygg og deploy
make build                        # Bygg app Docker-image
make build-migrate                # Bygg migrate Docker-image
make push                         # Push app-image til registry
make push-migrate                 # Push migrate-image til registry
make deploy                       # Deploy til k3s med Helm
make ship                         # Bygg, push og deploy i én kommando

# Utvikling
make up                           # Start docker-compose
make down                         # Stopp docker-compose
make migrate                      # Kjør migrasjoner lokalt
make test                         # Kjør alle tester
make vet                          # Kjør go vet
make clean                        # Slett kompilerte filer
```

## Helm

```bash
# Se hva som vil deployes
helm template task-mgr helm/ \
  --namespace task-mgr \
  --set image.tag=v0.0.1

# Deploy
helm upgrade --install task-mgr helm/ \
  --namespace task-mgr \
  --create-namespace \
  --set image.tag=v0.0.1

# Se installerte releases
helm list -n task-mgr

# Se historikk
helm history task-mgr -n task-mgr

# Rulle tilbake
helm rollback task-mgr -n task-mgr

# Slett release
helm uninstall task-mgr -n task-mgr
```

## kubectl

```bash
# Generelt
kubectl get all -n task-mgr                           # Se alle ressurser
kubectl get pods -n task-mgr                          # Se pods
kubectl get pods -n task-mgr -w                       # Watch pods (live)

# Feilsøking
kubectl describe pod -n task-mgr <pod-navn>           # Detaljer om pod
kubectl logs -n task-mgr <pod-navn>                   # Logger
kubectl logs -n task-mgr <pod-navn> -c <container>    # Logger fra spesifikk container
kubectl logs -n task-mgr <pod-navn> --previous        # Logger fra forrige instans

# Port-forward
kubectl port-forward -n task-mgr service/task-mgr-service 8080:80

# Exec
kubectl exec -it -n task-mgr <pod-navn> -- sh
kubectl exec -it -n task-mgr postgres-0 -- psql -U taskuser -d taskdb

# Rollout
kubectl rollout restart deployment/task-mgr -n task-mgr
kubectl rollout status deployment/task-mgr -n task-mgr
kubectl rollout history deployment/task-mgr -n task-mgr
```

## git (oh-my-zsh aliases)

```bash
gst                               # git status
gaa                               # git add --all
gcam "melding"                    # git commit -am (endrede filer)
gaa && gcam "melding"             # legg til nye filer og commit
glog                              # pen git log
gd                                # git diff
gp                                # git push
gl                                # git pull
gco -b <branch>                   # lag og bytt til ny branch
gm <branch>                       # merge branch
```

### Commit-konvensjon (Conventional Commits)

```
feat:     ny funksjonalitet
fix:      feilretting
docs:     dokumentasjonsendringer
refactor: omstrukturering uten funksjonsendring
test:     legge til eller endre tester
chore:    vedlikehold, avhengighetsoppdateringer
```

## Lokal JWT-testing

```bash
# Start mock JWKS-server (i egen terminal)
go run ./testdata/mock-jwks/

# Generer token
TOKEN=$(go run ./testdata/gen-token/)

# Test beskyttet endepunkt
curl -s -X POST localhost:8072/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title": "test"}' | jq .

# Test at åpent endepunkt ikke krever token
curl -s localhost:8072/tasks | jq .

# Test at manglende token gir 401
curl -si -X POST localhost:8072/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "test"}'
```

## curl (API-testing)

```bash
# Sett base URL
BASE=localhost:8072          # lokalt
BASE=homelab/task-mgr-api    # k3s

# Hent alle
curl -s $BASE/tasks | jq .

# Opprett (krever token)
curl -s -X POST $BASE/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title": "min task"}' | jq .

# Hent én
curl -s $BASE/tasks/1 | jq .

# Oppdater (krever token)
curl -s -X PATCH $BASE/tasks/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"done": true}' | jq .

# Slett (krever token)
curl -s -X DELETE $BASE/tasks/1 \
  -H "Authorization: Bearer $TOKEN"

# Helsesjekk
curl -s $BASE/healthz | jq .
```
