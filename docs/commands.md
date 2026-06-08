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

## kubectl

```bash
# Generelt
kubectl get all -n task-mgr                           # Se alle ressurser
kubectl get pods -n task-mgr                          # Se pods
kubectl get pods -n task-mgr -w                       # Watch pods (live)

# Deploy
kubectl apply -f k8s/                                 # Apply alle manifester
kubectl apply -f k8s/deployment.yaml                  # Apply én fil
kubectl delete -f k8s/                                # Slett alle ressurser

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

# Namespace
kubectl get all -n task-mgr                           # Alt i namespace
kubectl delete namespace task-mgr                     # Slett alt (forsiktig!)
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

Eksempler:
```bash
gcam "feat: REST API med alle CRUD endepunkter"
gcam "fix: håndter feil ved ugyldig task ID"
gcam "docs: legg til API-dokumentasjon"
gcam "refactor: flytt UpdateTaskRequest til handler-pakken"
```

## curl (API-testing)

```bash
# Lokalt
BASE=localhost:8072

# k3s
BASE=<domene>/task-mgr-api

# Hent alle
curl -s $BASE/tasks | jq .

# Opprett
curl -s -X POST $BASE/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "min task"}' | jq .

# Hent én
curl -s $BASE/tasks/1 | jq .

# Oppdater
curl -s -X PATCH $BASE/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"done": true}' | jq .

# Slett
curl -s -X DELETE $BASE/tasks/1
```
