// Entry point for the say-hi portfolio on Azure. Deployed at resource-group
// scope by .github/workflows/deploy.yml; see docs/azure/bicep.md.
targetScope = 'resourceGroup'

@description('Azure region. Defaults to the resource group location.')
param location string = resourceGroup().location

@description('Name of the container app; also used as a prefix for related resources.')
param appName string = 'say-hi'

@description('Container image to run, e.g. ghcr.io/sayhilel/say-hi:<git-sha>.')
param image string

@description('Version string shown by the `cloud` command.')
param appVersion string = 'manual'

@description('Enable the Cosmos DB free tier. Set false only if another account in the subscription already uses it.')
param cosmosFreeTier bool = true

@description('Optional Entra object ID of the site owner, granted read-only access to Cosmos DB data.')
param ownerPrincipalId string = ''

@description('Email for budget alerts. Leave empty to skip creating the budget.')
param budgetContactEmail string = ''

@description('Monthly budget amount in the billing currency.')
param budgetAmount int = 1

@description('Budget start date; must be the first day of a month.')
param budgetStartDate string = utcNow('yyyy-MM-01')

@description('Comma-separated custom hostnames already set up per docs/deploy.md, e.g. "sahilsinha.me,www.sahilsinha.me".')
param customDomains string = ''

var suffix = uniqueString(resourceGroup().id)
var cosmosDatabase = 'sayhi'
var tags = {
  app: appName
  managedBy: 'bicep'
  costCenter: 'portfolio'
}

module logs 'modules/logs.bicep' = {
  name: 'logs'
  params: {
    tags: tags
    location: location
    name: 'log-${appName}-${suffix}'
  }
}

module identity 'modules/identity.bicep' = {
  name: 'identity'
  params: {
    tags: tags
    location: location
    name: 'id-${appName}'
  }
}

module cosmos 'modules/cosmos.bicep' = {
  name: 'cosmos'
  params: {
    tags: tags
    location: location
    accountName: 'cosmos-${appName}-${suffix}'
    databaseName: cosmosDatabase
    appPrincipalId: identity.outputs.principalId
    enableFreeTier: cosmosFreeTier
    ownerPrincipalId: ownerPrincipalId
  }
}

module app 'modules/containerapp.bicep' = {
  name: 'containerapp'
  params: {
    tags: tags
    location: location
    appName: appName
    image: image
    appVersion: appVersion
    logAnalyticsWorkspaceName: logs.outputs.name
    identityId: identity.outputs.id
    identityClientId: identity.outputs.clientId
    cosmosEndpoint: cosmos.outputs.endpoint
    cosmosDatabase: cosmosDatabase
    customDomains: filter(map(split(customDomains, ','), d => trim(d)), d => !empty(d))
  }
}

module workbook 'modules/workbook.bicep' = {
  name: 'workbook'
  params: {
    tags: tags
    location: location
    workspaceId: logs.outputs.id
    appName: appName
  }
}

module budget 'modules/budget.bicep' = if (!empty(budgetContactEmail)) {
  name: 'budget'
  params: {
    name: 'budget-${appName}'
    amount: budgetAmount
    contactEmail: budgetContactEmail
    startDate: budgetStartDate
  }
}

output appUrl string = 'https://${app.outputs.fqdn}'
output latestRevision string = app.outputs.latestRevision
output cosmosAccount string = cosmos.outputs.accountName
