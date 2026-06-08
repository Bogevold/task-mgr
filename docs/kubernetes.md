# Kubernetes / k3s

## Oversikt over manifester

```
k8s/
  namespace.yaml        # Isolerer alle ressurser under task-mgr
  secret.yaml           # Sensitive verdier (ikke i git!)
  secrets.example.yaml  # Eksempel med dummy-verdier (i git)
  configmap.yaml        # Ikke-sensitiv konfigurasjon
  pg-pvc.yaml           # Persistent lagring for PostgreSQL
  pg-deployment.yaml    # PostgreSQL StatefulSet
  pg-service.yaml       # Eksponerer PostgreSQL internt i clusteret
  deployment.yaml       # App Deployment med init-containere
  service.yaml          # Eksponerer appen internt i clusteret
  middleware.yaml        # Traefik strip-prefix middleware
  ingress.yaml          # Eksponerer appen eksternt via Traefik
```

## Deploy fra scratch

Kjør i denne rekkefølgen:

```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/pg-pvc.yaml
kubectl apply -f k8s/pg-deployment.yaml
kubectl apply -f k8s/pg-service.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/middleware.yaml
kubectl apply -f k8s/ingress.yaml
```

Eller alle på én gang (rekkefølgen styres av Kubernetes selv):

```bash
kubectl apply -f k8s/
```

## Ressurser forklart

### Namespace

Isolerer alle task-mgr ressurser fra andre applikasjoner i clusteret.
Gjør det enkelt å slette alt på én gang: `kubectl delete namespace task-mgr`.

### Secret

Inneholder `POSTGRES_USER` og `POSTGRES_PASSWORD`.
`secret.yaml` skal aldri committes til git — legg den i `.gitignore`.
Bruk `secrets.example.yaml` som mal.

### StatefulSet (PostgreSQL)

`StatefulSet` fremfor `Deployment` for PostgreSQL fordi:
- Stabil nettverksidentitet (podnavnet endres ikke ved restart)
- Garantert rekkefølge ved oppstart og nedstengning
- Viktig for databaser som forventer stabil identitet

### PersistentVolumeClaim

Reserverer 1GB diskplass fra k3s sin lokale storage.
Data overlever pod-restart og redeployments.
`subPath: postgres-db-data` unngår problemer med `lost+found`-mapper i k3s.

### Deployment (App)

Kjører 2 replikaer for redundans. Inneholder to init-containere:

1. `vent-pa-postgres` — poller `pg_isready` til PostgreSQL svarer
2. `kjor-migrasjoner` — kjører `migrate up` mot databasen

App-containeren starter ikke før begge init-containere er fullført.

### Traefik Ingress

Eksponerer appen på `http://<domene>/task-mgr-api`.
`Middleware` stripper `/task-mgr-api`-prefikset før requesten når appen,
slik at appen ikke trenger å vite noe om URL-strukturen utenfra.

## Feilsøking

```bash
# Se status på alle ressurser
kubectl get all -n task-mgr

# Beskriv en pod (se events og konfigurasjon)
kubectl describe pod -n task-mgr <pod-navn>

# Se logger fra en container
kubectl logs -n task-mgr <pod-navn>

# Se logger fra en init-container
kubectl logs -n task-mgr <pod-navn> -c vent-pa-postgres
kubectl logs -n task-mgr <pod-navn> -c kjor-migrasjoner

# Port-forward for lokal testing
kubectl port-forward -n task-mgr service/task-mgr-service 8080:80

# Koble til PostgreSQL-poden direkte
kubectl exec -it -n task-mgr postgres-0 -- psql -U taskuser -d taskdb
```

## Oppdatere applikasjonen

```bash
# Bygg og push nytt image
docker build -t homelab:30500/task-mgr:latest .
docker push homelab:30500/task-mgr:latest

# Tving Kubernetes til å hente nytt image
kubectl rollout restart deployment/task-mgr -n task-mgr

# Følg rollout-status
kubectl rollout status deployment/task-mgr -n task-mgr
```
