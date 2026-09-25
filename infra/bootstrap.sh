#!/usr/bin/env bash
# One-time setup that lets GitHub Actions deploy to Azure without any stored
# secret. Run it once, signed in with `az login`, from the repo root:
#
#   ./infra/bootstrap.sh
#
# What it does:
#   1. Registers the resource providers the templates use.
#   2. Creates the resource group.
#   3. Creates an Entra ID app registration + service principal for GitHub.
#   4. Adds a federated credential that trusts ONLY this repository's
#      `production` GitHub environment (OIDC, no client secret).
#   5. Grants that principal Contributor on the resource group only.
#   6. Prints the GitHub repository variables to set.
#
# Safe to re-run: every step checks for existing resources first.
set -euo pipefail

REPO="${REPO:-sayhilel/say-hi}"
RG="${RG:-rg-say-hi}"
LOCATION="${LOCATION:-eastus}"
GH_ENVIRONMENT="${GH_ENVIRONMENT:-production}"
APP_DISPLAY_NAME="${APP_DISPLAY_NAME:-github-say-hi-deployer}"

echo "==> Using subscription"
SUBSCRIPTION_ID=$(az account show --query id -o tsv)
TENANT_ID=$(az account show --query tenantId -o tsv)
az account show --query "{name:name, id:id}" -o table

echo "==> Registering resource providers"
for ns in Microsoft.App Microsoft.OperationalInsights Microsoft.DocumentDB \
          Microsoft.ManagedIdentity Microsoft.Insights Microsoft.Consumption; do
  az provider register --namespace "$ns" --wait -o none
done

echo "==> Creating resource group $RG in $LOCATION"
az group create --name "$RG" --location "$LOCATION" \
  --tags app=say-hi managedBy=bicep costCenter=portfolio -o none
RG_ID=$(az group show --name "$RG" --query id -o tsv)

echo "==> Creating app registration $APP_DISPLAY_NAME"
CLIENT_ID=$(az ad app list --display-name "$APP_DISPLAY_NAME" --query "[0].appId" -o tsv)
if [[ -z "$CLIENT_ID" ]]; then
  CLIENT_ID=$(az ad app create --display-name "$APP_DISPLAY_NAME" --query appId -o tsv)
fi
SP_OBJECT_ID=$(az ad sp show --id "$CLIENT_ID" --query id -o tsv 2>/dev/null || true)
if [[ -z "$SP_OBJECT_ID" ]]; then
  SP_OBJECT_ID=$(az ad sp create --id "$CLIENT_ID" --query id -o tsv)
fi

echo "==> Adding federated credential for repo:$REPO:environment:$GH_ENVIRONMENT"
SUBJECT="repo:${REPO}:environment:${GH_ENVIRONMENT}"
EXISTING=$(az ad app federated-credential list --id "$CLIENT_ID" \
  --query "[?subject=='$SUBJECT'] | length(@)" -o tsv)
if [[ "$EXISTING" == "0" ]]; then
  az ad app federated-credential create --id "$CLIENT_ID" -o none --parameters "{
    \"name\": \"github-${GH_ENVIRONMENT}\",
    \"issuer\": \"https://token.actions.githubusercontent.com\",
    \"subject\": \"${SUBJECT}\",
    \"audiences\": [\"api://AzureADTokenExchange\"]
  }"
fi

echo "==> Granting Contributor on $RG (and nothing else)"
# Entra replication can lag a few seconds behind service principal creation.
for attempt in 1 2 3 4 5; do
  if az role assignment create --assignee-object-id "$SP_OBJECT_ID" \
       --assignee-principal-type ServicePrincipal \
       --role Contributor --scope "$RG_ID" -o none 2>/dev/null; then
    break
  fi
  if [[ "$attempt" == 5 ]]; then
    echo "ERROR: could not assign Contributor to $SP_OBJECT_ID on $RG_ID" >&2
    exit 1
  fi
  echo "    waiting for the service principal to replicate ($attempt/5)..."
  sleep 10
done

OWNER_PRINCIPAL_ID=$(az ad signed-in-user show --query id -o tsv 2>/dev/null || true)

cat <<VARS

Done. In GitHub: Settings -> Environments -> New environment "$GH_ENVIRONMENT",
then Settings -> Secrets and variables -> Actions -> Variables, and add:

  AZURE_CLIENT_ID        $CLIENT_ID
  AZURE_TENANT_ID        $TENANT_ID
  AZURE_SUBSCRIPTION_ID  $SUBSCRIPTION_ID
  AZURE_RESOURCE_GROUP   $RG
  OWNER_PRINCIPAL_ID     $OWNER_PRINCIPAL_ID   (optional: lets you read messages in Data Explorer)
  BUDGET_EMAIL           <your email>          (optional: creates the \$1 budget alert)

None of these are secrets: the IDs only identify the tenant, subscription and
app, and Azure accepts a token only when GitHub's OIDC token for this repo's
"$GH_ENVIRONMENT" environment is presented.
VARS
