@description('Tags applied to every resource.')
param tags object

@description('Azure region for the environment and app.')
param location string

@description('Container app name.')
param appName string

@description('Container image, e.g. ghcr.io/owner/say-hi:<sha>.')
param image string

@description('Version string shown by the `cloud` command (usually the git SHA).')
param appVersion string

@description('Log Analytics workspace that receives console and system logs.')
param logAnalyticsWorkspaceName string

@description('Resource ID of the user-assigned managed identity.')
param identityId string

@description('Client ID of the managed identity; DefaultAzureCredential reads it from AZURE_CLIENT_ID.')
param identityClientId string

@description('Cosmos DB endpoint.')
param cosmosEndpoint string

@description('Cosmos DB database name.')
param cosmosDatabase string

@description('Custom hostnames bound to the app, e.g. [\'sahilsinha.me\', \'www.sahilsinha.me\']. Each needs a managed certificate named cert-<hostname with dots replaced by dashes> in the environment; see docs/deploy.md.')
param customDomains array = []

resource workspace 'Microsoft.OperationalInsights/workspaces@2023-09-01' existing = {
  name: logAnalyticsWorkspaceName
}

// Consumption-only environment: there is no charge for the environment itself,
// and the first 180,000 vCPU-seconds, 360,000 GiB-seconds and 2 million
// requests each month are free.
resource environment 'Microsoft.App/managedEnvironments@2024-03-01' = {
  name: '${appName}-env'
  location: location
  tags: tags
  properties: {
    appLogsConfiguration: {
      destination: 'log-analytics'
      logAnalyticsConfiguration: {
        customerId: workspace.properties.customerId
        sharedKey: workspace.listKeys().primarySharedKey
      }
    }
    zoneRedundant: false
  }
}

resource app 'Microsoft.App/containerApps@2024-03-01' = {
  name: appName
  location: location
  tags: tags
  identity: {
    type: 'UserAssigned'
    userAssignedIdentities: {
      '${identityId}': {}
    }
  }
  properties: {
    managedEnvironmentId: environment.id
    configuration: {
      activeRevisionsMode: 'Single'
      ingress: {
        external: true
        targetPort: 8080
        transport: 'auto'
        allowInsecure: false // HTTP is redirected to HTTPS
        // Managed certificates are issued once by the CLI (they need the
        // hostname to exist first), then referenced here by their fixed name so
        // every redeploy keeps the bindings.
        customDomains: [
          for host in customDomains: {
            name: host
            certificateId: '${environment.id}/managedCertificates/cert-${replace(host, '.', '-')}'
            bindingType: 'SniEnabled'
          }
        ]
      }
    }
    template: {
      containers: [
        {
          name: appName
          image: image
          resources: {
            cpu: json('0.25')
            memory: '0.5Gi'
          }
          env: [
            { name: 'AZURE_CLIENT_ID', value: identityClientId }
            { name: 'COSMOS_ENDPOINT', value: cosmosEndpoint }
            { name: 'COSMOS_DATABASE', value: cosmosDatabase }
            { name: 'AZURE_REGION', value: location }
            { name: 'APP_VERSION', value: appVersion }
          ]
          // Probes use /healthz, not /readyz: if Cosmos DB has a blip the
          // portfolio should keep serving pages rather than be pulled out of
          // rotation. /readyz is used by the CI smoke test instead.
          probes: [
            {
              type: 'Startup'
              httpGet: { path: '/healthz', port: 8080 }
              periodSeconds: 2
              failureThreshold: 15
            }
            {
              type: 'Liveness'
              httpGet: { path: '/healthz', port: 8080 }
              periodSeconds: 30
            }
            {
              type: 'Readiness'
              httpGet: { path: '/healthz', port: 8080 }
              periodSeconds: 10
            }
          ]
        }
      ]
      scale: {
        // Scale to zero when idle so an unvisited portfolio costs nothing.
        minReplicas: 0
        maxReplicas: 1
        rules: [
          {
            name: 'http'
            http: {
              metadata: {
                concurrentRequests: '50'
              }
            }
          }
        ]
      }
    }
  }
}

output fqdn string = app.properties.configuration.ingress.fqdn
output environmentName string = environment.name
output latestRevision string = app.properties.latestRevisionName
