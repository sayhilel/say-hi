@description('Tags applied to every resource.')
param tags object

@description('Azure region for the workspace.')
param location string

@description('Log Analytics workspace name.')
param name string

@description('Daily ingestion cap in GB. 5 GB/month is free, so ~0.15 GB/day keeps the bill at zero.')
param dailyQuotaGb string = '0.15'

resource workspace 'Microsoft.OperationalInsights/workspaces@2023-09-01' = {
  name: name
  location: location
  tags: tags
  properties: {
    sku: {
      name: 'PerGB2018'
    }
    retentionInDays: 30 // 31 days of retention are included at no cost
    workspaceCapping: {
      dailyQuotaGb: json(dailyQuotaGb)
    }
  }
}

output id string = workspace.id
output name string = workspace.name
