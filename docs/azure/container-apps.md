# Azure Container Apps

**What it is:** a serverless container platform built on Kubernetes, KEDA (event-driven autoscaling), Envoy (ingress) and Dapr, with the cluster managed for you. You bring a container image. Azure provides HTTPS ingress, autoscaling (including to zero), revisions and log shipping. It's the Azure equivalent of Google Cloud Run, which is where this site used to run.

**Where:** [`infra/modules/containerapp.bicep`](../../infra/modules/containerapp.bicep)

## How it's used here

| Setting | Value | Why |
|---|---|---|
| Plan | Consumption-only environment | No charge for the environment itself. The monthly free grant covers the app |
| CPU / memory | 0.25 vCPU / 0.5 GiB | Smallest valid pair. Go + Fiber idles at ~15 MB |
| Scale | min 0, max 1, HTTP rule at 50 concurrent requests | Scale to zero means zero usage when nobody is visiting. Max 1 caps spend |
| Ingress | External, target port 8080, `allowInsecure: false` | Envoy terminates TLS with a platform certificate. HTTP is redirected to HTTPS |
| Revisions | Single-revision mode | Each deploy creates a new revision. Traffic moves only once it's healthy, which gives zero-downtime rollouts |
| Identity | User-assigned managed identity | Used to reach Cosmos DB with no secrets. See [managed-identity.md](managed-identity.md) |
| Probes | Startup, liveness and readiness on `/healthz` | The platform restarts a hung replica and holds traffic until startup finishes |
| Logs | `appLogsConfiguration` → Log Analytics | stdout/stderr go to `ContainerAppConsoleLogs_CL` and platform events to `ContainerAppSystemLogs_CL` |

### Why the probes use `/healthz` and not `/readyz`
`/readyz` checks Cosmos DB. If readiness depended on it, a short database outage would take the *whole site* out of rotation, even though Cosmos only backs two small features. So the probes check that the process is alive. `/readyz` is used by the deploy pipeline's smoke test, where a failure *should* block the release.

### App-side changes for the platform
- **Graceful shutdown** (`cmd/main.go`): Container Apps sends `SIGTERM` when scaling in or replacing a revision. The server drains in-flight requests for up to 10 s.
- **Client IP**: Envoy sets `X-Forwarded-For`. Fiber is configured with `ProxyHeader` so rate limiting and logs see the real client.
- **Runtime metadata**: the platform injects `CONTAINER_APP_NAME`, `CONTAINER_APP_REVISION` and `CONTAINER_APP_REPLICA_NAME`. The `cloud` command displays them (`internal/runtimeinfo`).

## Free-tier limits
Per subscription per month: 180,000 vCPU-seconds, 360,000 GiB-seconds, 2 million requests.
One replica at 0.25 vCPU running *non-stop* would use about 648,000 vCPU-s/month. That's why `minReplicas: 0` matters: the site only uses compute while someone is looking at it.

## Trade-off: cold starts
With scale-to-zero, the first visitor after an idle period waits for a replica to start (roughly 2–5 s for this small distroless image). The workbook's *Cold starts* tile counts these. Setting `minReplicas: 1` removes the delay but leaves the free grant after about 8 days of the month.

## In production I would...
- Use a workload-profiles environment with a **VNet** and **private endpoints** to Cosmos DB, and disable public network access on the database.
- Put **Azure Front Door + WAF** in front for global routing, DDoS protection and OWASP rules.
- Run with `minReplicas >= 2` across **availability zones**, and use multi-revision mode for **blue/green or canary** traffic splitting.
