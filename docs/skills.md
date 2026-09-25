# Skills demonstrated

A map from each skill to where it shows up in this repo, and why it matters to a regulated, security-conscious employer such as a financial institution.

| Skill | Where to look | What it shows / why it matters |
|---|---|---|
| **Infrastructure as Code (Bicep)** | `infra/` · [bicep.md](azure/bicep.md) | The whole environment is reproducible from git, peer-reviewed like code, and every change has a deployment record. That makes it *auditable*, which regulators expect for production change management |
| **Change preview** | `what-if` step in `deploy.yml` | Every infrastructure change is previewed before it's applied, which reduces the risk of accidental deletions |
| **Passwordless CI/CD (OIDC federation)** | `deploy.yml`, `bootstrap.sh` · [github-actions-oidc.md](azure/github-actions-oidc.md) | There are no long-lived cloud credentials in GitHub. Removing credentials that can leak is a direct response to a common cause of cloud breaches |
| **Managed identity + data-plane RBAC** | `identity.bicep`, `cosmos.bicep`, `internal/store/cosmos.go` · [managed-identity.md](azure/managed-identity.md) | The app has no connection strings, and database keys are disabled at the account level. It's a zero-secret architecture |
| **Least privilege** | Deployer: Contributor on one RG. App: data contributor on one database. Owner: read-only | Each identity can do exactly its job and nothing more |
| **Container security** | `Dockerfile` · [registry.md](azure/registry.md) | Multi-stage build, distroless, non-root, static binary. Smaller attack surface and fewer CVEs |
| **Web security basics** | `cmd/main.go`, `handlers.SubmitContact` | Security headers (HSTS, nosniff, frame options, referrer policy) via `helmet`; HTTPS-only ingress; input validation with length limits; honeypot; per-IP rate limiting; `HttpOnly`/`SameSite` cookies; panic recovery |
| **Serverless containers / autoscaling** | `containerapp.bicep` · [container-apps.md](azure/container-apps.md) | Scale-to-zero, revisions, zero-downtime rollouts, health probes, graceful SIGTERM shutdown |
| **NoSQL data modelling** | `cosmos.bicep`, `internal/store` · [cosmos-db.md](azure/cosmos-db.md) | Partition-key choice, RU budgeting, server-side atomic `increment` patch instead of read-modify-write, idempotent document creation |
| **Resilience** | Timeouts on every external call, in-memory fallback, probes that don't depend on the database | A failing dependency degrades one feature rather than taking the whole site down |
| **Observability** | `log/slog` JSON logs, `workbook.bicep` · [log-analytics.md](azure/log-analytics.md) | Structured logs, KQL, latency percentiles, error tracking, and a dashboard that is itself code |
| **FinOps / cost engineering** | Throughput cap, log cap, scale-to-zero, budget, tags · [cost-management.md](azure/cost-management.md) | Runs at $0 by design, prevents spend rather than only alerting on it, and tags every resource for chargeback |
| **Cloud migration** | GCP Cloud Run → Azure Container Apps | Mapped managed services one to one and kept the app portable. The same container runs locally, on Cloud Run or on Azure |
| **Testing & quality gates** | `internal/**/_test.go`, `ci.yml` | Handler tests for routing, validation, the visitor counter and bounds checks; `go vet`, `gofmt`, race detector and `govulncheck` (dependency CVE scan) on every PR |
| **12-factor configuration** | `store.FromEnv`, `runtimeinfo` | Config comes from the environment. The same image is promoted unchanged, and local dev needs no cloud access |
| **Technical writing** | `docs/` | Decisions, trade-offs and "in production I would..." sections, with the reasoning behind each choice |

## Things I would add in a production setting
These are documented in each service page and left out here only because they aren't free:
- private networking (VNet, private endpoints)
- Front Door + WAF
- Key Vault for any third-party secrets
- Application Insights tracing and alert rules
- Defender for Cloud
- image signing and SBOMs
- Azure Policy guardrails
- multi-region disaster recovery
