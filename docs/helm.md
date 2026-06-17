# Helm

## Oversikt

Applikasjonen deployes til k3s via et Helm chart. Helm håndterer templating av Kubernetes-manifester
og gjør det enkelt å deploye til forskjellige miljøer med ulike verdier.

## Chart-struktur

```
helm/
  Chart.yaml              # Chart-metadata og versjon
  values.yaml             # Standardverdier
  templates/
    _helpers.tpl          # Gjenbrukbare hjelpefunksjoner
    deployment.yaml       # App Deployment med init-containere
    service.yaml          # App Service
    ingress.yaml          # Traefik Ingress
    middleware.yaml       # Traefik strip-prefix Middleware
    postgres/
      statefulset.yaml    # PostgreSQL StatefulSet
      service.yaml        # PostgreSQL Service
      pvc.yaml            # PersistentVolumeClaim
```

## values.yaml

```yaml
registry: homelab:30500

app:
  port: 8072
  replicas: 2
  env:
    jwksurl: "https://gitlab.example.com/-/jwks"
    jwtaud: "https://din-app.com"
    alwdns: "gruppe-a,gruppe-b"

postgres:
  database: taskdb
  port: 5432
  storage: 1Gi

ingress:
  host: homelab
  path: /task-mgr-api

image:
  app: task-mgr
  migrate: task-mgr-migrate
```

`image.tag` settes alltid av Makefilen ved deploy — aldri i `values.yaml`.

## Hjelpefunksjoner (_helpers.tpl)

| Funksjon | Returnerer |
|----------|-----------|
| `app.name` | Chart-navn (task-mgr) |
| `app.fullname` | Release-navn + chart-navn |
| `app.secretName` | `<fullname>-secret` |
| `app.labels` | Standard Kubernetes-labels |
| `app.selectorLabels` | Selector-labels for pods |

## Deploy

```bash
# Første deploy
kubectl apply -f k8s/secret.yaml
helm upgrade --install task-mgr helm/ \
  --namespace task-mgr \
  --create-namespace \
  --set image.tag=v0.0.1

# Oppdater
helm upgrade task-mgr helm/ \
  --namespace task-mgr \
  --set image.tag=v0.1.0

# Via Makefile
make ship
```

## Nyttige Helm-kommandoer

```bash
# Se hva som vil deployes uten å deploye
helm template task-mgr helm/ \
  --namespace task-mgr \
  --set image.tag=v0.0.1

# Se installerte releases
helm list -n task-mgr

# Se historikk
helm history task-mgr -n task-mgr

# Rulle tilbake til forrige versjon
helm rollback task-mgr -n task-mgr

# Slett release (beholder PVC)
helm uninstall task-mgr -n task-mgr
```

## Secrets

Secrets håndteres ikke av Helm — de opprettes manuelt:

```bash
# Kopier eksempel og fyll inn ekte verdier
cp k8s/secrets.example.yaml k8s/secret.yaml
kubectl apply -f k8s/secret.yaml
```

`k8s/secret.yaml` skal aldri committes til git (ligger i `.gitignore`).
