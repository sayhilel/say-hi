# GitHub Actions + OIDC workload identity federation

**What it is:** GitHub Actions runs the CI/CD. **OpenID Connect (OIDC) federation** lets a workflow sign in to Azure *without any stored secret*:
1. GitHub issues a short-lived signed token that says "this job is running in repo X, environment Y".
2. Microsoft Entra ID trusts that token for one specific app registration and exchanges it for an Azure access token.

**Where:** [`.github/workflows/ci.yml`](../../.github/workflows/ci.yml), [`.github/workflows/deploy.yml`](../../.github/workflows/deploy.yml), [`infra/bootstrap.sh`](../../infra/bootstrap.sh)

## Pipeline

```
push to main
   │
   ├─ ci (reusable workflow, also runs on every PR)
   │    ├─ go:     gofmt check · go vet · go test -race · govulncheck
   │    ├─ bicep:  az bicep lint · build · build-params
   │    └─ docker: build image (cached with GitHub Actions cache)
   │
   ├─ build   → push ghcr.io/sayhilel/say-hi:<git sha> (+ :latest)
   │
   └─ deploy  (GitHub environment: production)
        ├─ azure/login via OIDC
        ├─ az deployment group what-if     ← reviewable plan in the log
        ├─ az deployment group create      ← new Container Apps revision
        └─ smoke test: /healthz (waits out cold start) → /readyz (Cosmos via MI)
                       → `cloud` command shows the deployed commit SHA
```

Images are tagged with the **commit SHA**, so every running revision can be traced to exact source code. Rolling back means redeploying an earlier SHA.

## How the trust is set up
`infra/bootstrap.sh` (run once by a human with `az login`):
1. Creates an Entra **app registration** and service principal: `github-say-hi-deployer`.
2. Adds a **federated credential**:
   - issuer: `https://token.actions.githubusercontent.com`
   - subject: `repo:sayhilel/say-hi:environment:production`
   - audience: `api://AzureADTokenExchange`

   Only jobs that run in this repo's `production` environment can get a token. Forks, other branches' PR jobs and other repos can't.
3. Grants the principal **Contributor on the `rg-say-hi` resource group only**, not the subscription. Owner or User Access Administrator isn't needed, because the only role assignments the templates create are Cosmos DB *data-plane* assignments, which are Cosmos resources rather than Azure RBAC.

The workflow reads `AZURE_CLIENT_ID`, `AZURE_TENANT_ID` and `AZURE_SUBSCRIPTION_ID` from repository **variables**, not secrets. They're identifiers, not credentials. The workflow's `permissions:` block grants `id-token: write` only to the deploy job, and `packages: write` only to the build job.

## Why this matters
A long-lived client secret in CI is one of the most common ways cloud credentials leak: it gets copied, logged, never rotated and outlives the people who created it. With OIDC:
- there's nothing to rotate or leak;
- every token lasts minutes and is bound to a specific repo and environment;
- sign-ins show up in Entra sign-in logs with the workflow run's claims.

## In production I would...
- Add **required reviewers** and **branch protection** to the `production` environment for a four-eyes approval before deploys (a checkbox in GitHub settings).
- Sign images (Sigstore/cosign or Notation), generate an **SBOM**, and verify the signature at deploy time.
- Pin third-party actions to commit SHAs and use Dependabot to update them.
