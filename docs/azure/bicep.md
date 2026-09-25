# Bicep (Infrastructure as Code)

**What it is:** Bicep is Azure's own declarative IaC language. It compiles to ARM (Azure Resource Manager) JSON templates. Compared with Terraform it has no state file (Azure itself is the state), day-one support for new resource types, and type-checked resource properties.

**Where:** [`infra/`](../../infra)

```
infra/
├── main.bicep           # composes the modules; resource-group scope
├── main.bicepparam      # typed parameter file
├── bootstrap.sh         # one-time identity/RG setup (see github-actions-oidc.md)
└── modules/
    ├── logs.bicep          Log Analytics workspace with a daily cap
    ├── identity.bicep      user-assigned managed identity
    ├── cosmos.bicep        free-tier account, database, containers, data-plane RBAC
    ├── containerapp.bicep  environment + app, probes, scaling, ingress
    ├── workbook.bicep      operations dashboard generated from a list of KQL tiles
    └── budget.bicep        $1 monthly budget with email alerts
```

## Patterns used
- **Modules with explicit inputs and outputs.** For example, `identity` outputs `principalId`, which `cosmos` needs for its role assignment and `containerapp` needs to attach the identity. ARM works out the dependency order from these references, so there's no hand-written `dependsOn` except for child resources that need it.
- **Deterministic unique names.** `uniqueString(resourceGroup().id)` gives globally unique names (the Cosmos account needs one) that stay the same across redeploys.
- **Idempotent role assignments.** Role assignment names come from `guid(scope, principal, role)`, so redeploying never creates duplicates.
- **Conditional resources.** The budget, the owner's data-reader role and the custom domain binding are only created when their parameter is set.
- **Code-generated config.** `workbook.bicep` builds the dashboard JSON with `map()` over a list of `{title, query, visualization}`, so adding a tile is one line.
- **Consistent tags.** `app`, `managedBy` and `costCenter` are applied to every resource for cost reporting and ownership.
- **No secrets in parameters.** The only sensitive value, the Log Analytics shared key, is read with `listKeys()` inside the module. It is never a parameter or an output.

## Safe deployments
The pipeline runs **`az deployment group what-if`** before every deployment. It prints a plan of what will be created, changed or deleted, similar to `terraform plan`, into the job log for review. Deployments are named `say-hi-<run number>`, so the resource group's *Deployments* blade keeps an audit trail of every infrastructure change and its inputs.

## Running it by hand
```bash
az deployment group what-if -g rg-say-hi --parameters infra/main.bicepparam
az deployment group create  -g rg-say-hi --parameters infra/main.bicepparam \
  --parameters image=ghcr.io/sayhilel/say-hi:latest
```

## In production I would...
- Use **deployment stacks** with deny settings, so resources can only be changed through the pipeline.
- Run a policy and security scanner (PSRule for Azure, Checkov) in CI, plus **Azure Policy** assignments for guardrails such as allowed regions, required tags and denying public endpoints.
- Add separate parameter files per environment (dev/test/prod) promoted through the same pipeline.
