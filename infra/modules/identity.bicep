@description('Tags applied to every resource.')
param tags object

@description('Azure region for the identity.')
param location string

@description('User-assigned managed identity name.')
param name string

// A user-assigned identity (rather than system-assigned) exists before the
// container app does, so Cosmos DB role assignments can be granted in the same
// deployment without a chicken-and-egg problem.
resource identity 'Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31' = {
  name: name
  location: location
  tags: tags
}

output id string = identity.id
output clientId string = identity.properties.clientId
output principalId string = identity.properties.principalId
