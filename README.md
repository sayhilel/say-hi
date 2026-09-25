# say-hi

A terminal-style portfolio, written in Go ([Fiber](https://gofiber.io)) and [HTMX](https://htmx.org), running on **Microsoft Azure** within the free tier.

You navigate by typing commands:

| Command | What it does |
|---|---|
| `whoami` | About me |
| `exp` | Professional experience |
| `showcase` | Projects (use ↑/↓) |
| `ping` | Contact form (stored in Cosmos DB) |
| `cloud` | How this site is deployed, with live revision/replica/region info |
| `more` | Open resume |
| `clear` | Clear the terminal |

## Azure at a glance
- **Compute:** Azure Container Apps (scale to zero, distroless non-root container)
- **Data:** Cosmos DB for NoSQL, free tier, keys disabled, accessed via **managed identity**
- **IaC:** Bicep modules in [`infra/`](infra)
- **CI/CD:** GitHub Actions with **OIDC federation**, so there are no stored cloud credentials. Each deploy runs `what-if`, then deploys, then smoke-tests
- **Observability:** structured JSON logs → Log Analytics (KQL) → Workbook dashboard
- **Cost:** $0/month by design, with hard caps and a budget alert

Architecture, per-service write-ups and skills: **[docs/](docs/README.md)**.
Deploying it yourself: **[docs/deploy.md](docs/deploy.md)**.

## Run locally
```bash
go run ./cmd        # http://localhost:8080, in-memory store, no Azure needed
go test ./...
```

## Layout
```
cmd/                  entry point: routes, middleware, graceful shutdown
internal/handlers     terminal commands, contact form, visitor counter
internal/store        Cosmos DB (managed identity) / in-memory storage
internal/runtimeinfo  Container Apps runtime metadata for `cloud`
internal/projects     `showcase` data from public/data/projects.toml
views/, public/       HTML templates and static assets
infra/                Bicep + one-time bootstrap script
.github/workflows/    CI and deploy pipelines
docs/                 architecture and service documentation
```
