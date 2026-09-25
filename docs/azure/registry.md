# Container registry: GHCR, and why not ACR

**Where:** the `build` job in [`deploy.yml`](../../.github/workflows/deploy.yml)

## Decision
Images are pushed to **GitHub Container Registry** (`ghcr.io/sayhilel/say-hi`) as a **public** package. Container Apps pulls public images anonymously, so the app configuration needs no registry credentials at all.

## Why not Azure Container Registry?
ACR is the natural choice on Azure and is what I would use in an enterprise. But even its Basic tier has a fixed daily charge (about $5/month) once the free-account credit ends, and this project's constraint is $0/month. Everything in the image is already public on GitHub, so a public registry exposes nothing new.

## What would change with ACR
The change is small, which is part of why the design works:
1. Add an `acr.bicep` module (`Microsoft.ContainerRegistry/registries`, `adminUserEnabled: false`).
2. Grant the app's managed identity the **AcrPull** role on the registry. Also grant the GitHub deployer principal **AcrPush**; this is a real Azure RBAC assignment, so the deployer then needs *Role Based Access Control Administrator*, limited to those roles.
3. In the container app set `configuration.registries: [{ server: '<acr>.azurecr.io', identity: <identity id> }]`, so pulls use the managed identity and there are still no passwords.
4. In the workflow, replace the GHCR login with `az acr login` (via the same OIDC session), or use `az acr build` to build inside Azure.

Enterprise ACR features worth mentioning: Premium-tier private endpoints, geo-replication, Defender for Containers vulnerability scanning, and content trust/Notation signing.

## Image hardening (applies either way)
The [`Dockerfile`](../../Dockerfile) is multi-stage:
- **Build stage:** `golang:1.25-alpine`, static binary (`CGO_ENABLED=0`, `-trimpath -ldflags="-s -w"`).
- **Runtime stage:** `gcr.io/distroless/static-debian12:nonroot`. It has no shell and no package manager, runs as UID 65532, and includes CA certificates for TLS to Cosmos DB. The final image contains only the binary, `views/` and `public/`.

Fewer packages means fewer CVEs to patch and very little an attacker can use if the process is ever compromised.
