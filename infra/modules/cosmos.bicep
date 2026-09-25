@description('Tags applied to every resource.')
param tags object

@description('Azure region for the Cosmos DB account.')
param location string

@description('Globally unique Cosmos DB account name.')
param accountName string

@description('SQL database name.')
param databaseName string

@description('Principal ID of the managed identity that the app runs as.')
param appPrincipalId string

@description('Optional Entra object ID granted read-only data access, so the owner can read contact messages in Data Explorer. Leave empty to skip.')
param ownerPrincipalId string = ''

@description('Enable the free tier (1000 RU/s + 25 GB free forever). Only one account per subscription may use it.')
param enableFreeTier bool = true

// Built-in data-plane roles.
var dataReaderRoleId = '00000000-0000-0000-0000-000000000001'
var dataContributorRoleId = '00000000-0000-0000-0000-000000000002'

resource account 'Microsoft.DocumentDB/databaseAccounts@2024-11-15' = {
  name: accountName
  location: location
  tags: tags
  kind: 'GlobalDocumentDB'
  properties: {
    databaseAccountOfferType: 'Standard'
    enableFreeTier: enableFreeTier
    locations: [
      {
        locationName: location
        failoverPriority: 0
        isZoneRedundant: false
      }
    ]
    consistencyPolicy: {
      defaultConsistencyLevel: 'Session'
    }
    // Hard cap on provisioned throughput so the account can never exceed the
    // free tier allowance, even if someone adds a container by hand.
    capacity: {
      totalThroughputLimit: 1000
    }
    // Keys are disabled: the only way in is Microsoft Entra ID (managed identity
    // for the app, your own account for Data Explorer).
    disableLocalAuth: true
    minimalTlsVersion: 'Tls12'
    publicNetworkAccess: 'Enabled'
  }
}

resource database 'Microsoft.DocumentDB/databaseAccounts/sqlDatabases@2024-11-15' = {
  parent: account
  name: databaseName
  properties: {
    resource: {
      id: databaseName
    }
    // Throughput shared by every container in the database; 1000 RU/s is the
    // full free tier allowance.
    options: {
      throughput: 1000
    }
  }
}

resource messages 'Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers@2024-11-15' = {
  parent: database
  name: 'messages'
  properties: {
    resource: {
      id: 'messages'
      partitionKey: {
        paths: ['/kind']
        kind: 'Hash'
      }
    }
  }
}

resource stats 'Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers@2024-11-15' = {
  parent: database
  name: 'stats'
  properties: {
    resource: {
      id: 'stats'
      partitionKey: {
        paths: ['/id']
        kind: 'Hash'
      }
    }
  }
}

// Data-plane RBAC, scoped to this one database rather than the whole account.
resource appDataAccess 'Microsoft.DocumentDB/databaseAccounts/sqlRoleAssignments@2024-11-15' = {
  parent: account
  name: guid(account.id, appPrincipalId, dataContributorRoleId)
  properties: {
    roleDefinitionId: '${account.id}/sqlRoleDefinitions/${dataContributorRoleId}'
    principalId: appPrincipalId
    scope: '${account.id}/dbs/${databaseName}'
  }
  dependsOn: [
    database
  ]
}

resource ownerDataAccess 'Microsoft.DocumentDB/databaseAccounts/sqlRoleAssignments@2024-11-15' = if (!empty(ownerPrincipalId)) {
  parent: account
  name: guid(account.id, ownerPrincipalId, dataReaderRoleId)
  properties: {
    roleDefinitionId: '${account.id}/sqlRoleDefinitions/${dataReaderRoleId}'
    principalId: ownerPrincipalId
    scope: '${account.id}/dbs/${databaseName}'
  }
  dependsOn: [
    database
  ]
}

output endpoint string = account.properties.documentEndpoint
output accountName string = account.name
