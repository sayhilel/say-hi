using 'main.bicep'

// The pipeline overrides `image` and `appVersion` on every deploy; these
// defaults are for a manual `az deployment group create`.
param image = 'ghcr.io/sayhilel/say-hi:latest'
param appVersion = 'manual'
// Budget alerts go to the BUDGET_EMAIL repository variable (see deploy.yml).
param budgetAmount = 1
