# Deploying say-hi to Azure

## Prerequisites
- An Azure account. The free account gives $200 of credit for 30 days plus always-free services; everything here fits in the always-free tier.
- [Azure CLI](https://learn.microsoft.com/cli/azure/install-azure-cli) 2.60 or newer, signed in with `az login`.
- Admin access to the GitHub repo.
- No other Cosmos DB account in the subscription already using the free tier (only one is allowed per subscription).

## 1. Bootstrap (once)
```bash
az login
./infra/bootstrap.sh          # override with REPO=..., RG=..., LOCATION=... if needed
```
It registers resource providers, creates `rg-say-hi`, creates the GitHub deployer identity with an OIDC federated credential, and grants it Contributor on the resource group. It prints the variables for the next step.

## 2. Configure GitHub (once)
1. **Settings → Environments → New environment → `production`.** Optionally add yourself as a required reviewer to get an approval gate.
2. **Settings → Secrets and variables → Actions → Variables**, add:

| Variable | Required | Value |
|---|---|---|
| `AZURE_CLIENT_ID` | yes | from bootstrap output |
| `AZURE_TENANT_ID` | yes | from bootstrap output |
| `AZURE_SUBSCRIPTION_ID` | yes | from bootstrap output |
| `AZURE_RESOURCE_GROUP` | yes | `rg-say-hi` |
| `BUDGET_EMAIL` | recommended | where budget alerts go |
| `OWNER_PRINCIPAL_ID` | optional | your Entra object ID; lets you read messages in Data Explorer |

## 3. First deploy
Push to `main`, or run the **deploy** workflow manually.

The first run fails at the deploy step with an image-pull error. This is expected: new GHCR packages are **private**. After the `build` job has pushed the image:
1. Go to GitHub → your profile → **Packages → say-hi → Package settings → Change visibility → Public**.
2. Re-run the failed `deploy` job.

The job summary shows the URL, `https://say-hi.<random>.<region>.azurecontainerapps.io`. Every later push to `main` deploys automatically.

## 4. Verify
- Open the URL and type `cloud`. You should see `Azure Container Apps`, the region, revision, replica, commit SHA, and `Azure Cosmos DB ...: connected`.
- Type `ping` and send a test message. It should appear in **Cosmos DB → Data Explorer → sayhi → messages**. This needs `OWNER_PRINCIPAL_ID` to have been set.
- **Monitor → Workbooks → say-hi operations** should show requests. Logs take about 2–5 minutes to arrive.
- **Cost Management → Cost analysis** for `rg-say-hi` should show $0.00.

## 5. Custom domain (optional, e.g. sahilsinha.me)
Container Apps issues free managed TLS certificates. The certificates are created once with the CLI, because a certificate can only be issued after the hostname is on the app. Bicep then refers to them by a fixed name, `cert-<hostname-with-dashes>`, so every deploy keeps the bindings.

**DNS.** Get the values:
```bash
az containerapp show -g rg-say-hi -n say-hi --query properties.configuration.ingress.fqdn -o tsv
az containerapp show -g rg-say-hi -n say-hi --query properties.customDomainVerificationId -o tsv
az containerapp env show -g rg-say-hi -n say-hi-env --query properties.staticIp -o tsv
```

Then create these records at your DNS provider:

| Record | Type | Value |
|---|---|---|
| `www` | CNAME | app FQDN |
| `@` (apex) | A | environment static IP (or CNAME/ALIAS if the provider flattens it) |
| `asuid.www` | TXT | verification ID |
| `asuid` | TXT | verification ID |

**Bind (once per hostname):**
```bash
RG=rg-say-hi; APP=say-hi; ENV=say-hi-env
bind() { # $1 = hostname, $2 = CNAME (subdomain) or HTTP (apex)
  local name="cert-${1//./-}"
  az containerapp hostname add -g $RG -n $APP --hostname "$1"
  az containerapp env certificate create -g $RG -n $ENV --hostname "$1" \
    --validation-method "$2" --certificate-name "$name"
  # wait until: az containerapp env certificate list -g $RG -n $ENV --managed-certificates-only
  #             shows provisioningState "Succeeded" for $name
  az containerapp hostname bind -g $RG -n $APP --hostname "$1" --environment $ENV \
    --certificate "$name"
}
bind www.sahilsinha.me CNAME
bind sahilsinha.me HTTP
```

**Keep them on redeploy:** set the repository variable `CUSTOM_DOMAINS=sahilsinha.me,www.sahilsinha.me`.

## Rollback
Each deploy creates a new revision tagged with the commit SHA. Either:
- re-run the **deploy** workflow for an earlier commit, or
- use the CLI immediately:
  ```bash
  az containerapp update -g rg-say-hi -n say-hi --image ghcr.io/sayhilel/say-hi:<older-sha>
  ```

## Teardown
```bash
az group delete --name rg-say-hi --yes
az ad app delete --id <AZURE_CLIENT_ID>
```

## Troubleshooting

| Symptom | Cause / fix |
|---|---|
| `AADSTS700213: No matching federated identity record` | The job isn't running in the `production` environment, or the repo name differs from the bootstrap `REPO` |
| Revision stuck *Activating*, `ImagePullBackOff` in system logs | The GHCR package is still private (see step 3) |
| `cloud` shows Cosmos `unreachable`, `/readyz` returns 503 | Role assignment still propagating (wait about 5 minutes), or `AZURE_CLIENT_ID` env mismatch. Check `ContainerAppConsoleLogs_CL` for the error |
| Cosmos deploy fails with *free tier already used* | Deploy with `--parameters cosmosFreeTier=false` (this costs money), or delete the other free-tier account |
| `The subscription is not registered to use namespace` | Re-run `bootstrap.sh`; it registers the providers |

## Running locally
```bash
go run ./cmd                 # in-memory store, http://localhost:8080
air                          # live reload (see .air.toml)

# Against the real Cosmos DB, using your az login:
COSMOS_ENDPOINT=https://<account>.documents.azure.com:443/ go run ./cmd
```
