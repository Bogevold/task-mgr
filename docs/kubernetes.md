# Kubernetes / k3s

## Oversikt over ressurser

Applikasjonen deployes via Helm. Helm genererer disse Kubernetes-ressursene:

```
Namespace: task-mgr
  Deployment: task-mgr         (2 replikaer)
    InitContainer: vent-pa-postgres
    InitContainer: kjor-migrasjoner
    Container: app
  Service: task-mgr-service    (ClusterIP, port 80→8072)
  Ingress: task-mgr-ingress    (Traefik, /task-mgr-api)
  Middleware: strip-task-mgr-api (Traefik strip-prefix)
  StatefulSet: postgres        (1 replika)
  Service: postgres            (ClusterIP, port 5432)
  PersistentVolumeClaim: postgres-pvc (1Gi)
```

Secrets opprettes manuelt utenfor Helm:
```
Secret: task-mgr-secret        (POSTGRES_USER, POSTGRES_PASSWORD)
```

## Første deploy

```bash
# Opprett namespace og secret
kubectl create namespace task-mgr
kubectl apply -f k8s/secret.yaml

# Deploy med Helm
helm upgrade --install task-mgr helm/ \
  --namespace task-mgr \
  --create-namespace \
  --set image.tag=v0.0.1
```

## Oppdatere applikasjonen

```bash
# Via Makefile (anbefalt)
make ship

# Manuelt
docker build -t homelab:30500/task-mgr:v0.1.0 .
docker push homelab:30500/task-mgr:v0.1.0
helm upgrade task-mgr helm/ \
  --namespace task-mgr \
  --set image.tag=v0.1.0
```

## Init-containere

Deployment har to init-containere som kjører i rekkefølge før app-containeren starter:

`vent-pa-postgres` — poller `pg_isready` til PostgreSQL svarer. Forhindrer at appen
starter før databasen er klar.

`kjor-migrasjoner` — kjører `migrate up` mot databasen. Sikrer at skjemaet alltid
er oppdatert før appen starter. Kjøres én gang per pod-oppstart.

## Health checks

`readinessProbe` — sjekker `/healthz` hvert 10. sekund. Hvis den feiler fjernes
poden fra load balanceren inntil databasen er tilgjengelig igjen.

`livenessProbe` — sjekker `/healthz` hvert 20. sekund. Hvis den feiler 3 ganger
på rad restarter Kubernetes poden.

## Traefik Ingress

Eksponerer appen på `http://homelab/task-mgr-api`.
`Middleware` stripper `/task-mgr-api`-prefikset før requesten når appen.

## PersistentVolumeClaim

`postgres-pvc` reserverer 1GB diskplass fra k3s sin lokale storage.
Data overlever pod-restart og redeployments.
`subPath: postgres-db-data` unngår problemer med `lost+found`-mapper i k3s.

Sletter du PVC-en mister du alle data — vær forsiktig med `kubectl delete pvc`.

## Feilsøking

```bash
# Se status på alle ressurser
kubectl get all -n task-mgr

# Se pods
kubectl get pods -n task-mgr -w

# Beskriv en pod (se events og konfigurasjon)
kubectl describe pod -n task-mgr <pod-navn>

# Se logger fra app-container
kubectl logs -n task-mgr <pod-navn>

# Se logger fra init-containere
kubectl logs -n task-mgr <pod-navn> -c vent-pa-postgres
kubectl logs -n task-mgr <pod-navn> -c kjor-migrasjoner

# Port-forward for lokal testing
kubectl port-forward -n task-mgr service/task-mgr-service 8080:80

# Koble til PostgreSQL-poden direkte
kubectl exec -it -n task-mgr postgres-0 -- psql -U taskuser -d taskdb

# Se Helm-historikk
helm history task-mgr -n task-mgr

# Rulle tilbake
helm rollback task-mgr -n task-mgr
```
