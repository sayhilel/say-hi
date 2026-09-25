@description('Tags applied to every resource.')
param tags object

@description('Azure region for the workbook.')
param location string

@description('Resource ID of the Log Analytics workspace the queries run against.')
param workspaceId string

@description('Container app name used to filter log rows.')
param appName string

// Every request is logged by the app as one JSON line on stdout; Container Apps
// ships it to ContainerAppConsoleLogs_CL with the line in Log_s.
var requests = 'ContainerAppConsoleLogs_CL | where ContainerAppName_s == "${appName}" | extend l = parse_json(Log_s) | where l.msg == "request"'

var tiles = [
  {
    title: 'Requests per hour'
    query: '${requests} | summarize requests = count() by bin(TimeGenerated, 1h)'
    visualization: 'timechart'
  }
  {
    title: 'Latency (ms): p50 / p95'
    query: '${requests} | summarize p50 = percentile(toint(l.durationMs), 50), p95 = percentile(toint(l.durationMs), 95) by bin(TimeGenerated, 1h)'
    visualization: 'linechart'
  }
  {
    title: 'Responses by status code'
    query: '${requests} | summarize count() by status = tostring(l.status)'
    visualization: 'piechart'
  }
  {
    title: 'Top paths'
    query: '${requests} | summarize hits = count() by path = tostring(l.path) | top 10 by hits'
    visualization: 'table'
  }
  {
    title: 'Cold starts (replica woke from zero)'
    query: 'ContainerAppConsoleLogs_CL | where ContainerAppName_s == "${appName}" | extend l = parse_json(Log_s) | where l.msg == "listening" | summarize coldStarts = count() by bin(TimeGenerated, 1d)'
    visualization: 'barchart'
  }
  {
    title: 'Contact messages received'
    query: 'ContainerAppConsoleLogs_CL | where ContainerAppName_s == "${appName}" | extend l = parse_json(Log_s) | where l.msg == "contact message saved" | summarize messages = count() by bin(TimeGenerated, 1d)'
    visualization: 'barchart'
  }
  {
    title: 'Recent errors'
    query: 'ContainerAppConsoleLogs_CL | where ContainerAppName_s == "${appName}" | extend l = parse_json(Log_s) | where l.level == "ERROR" | project TimeGenerated, msg = tostring(l.msg), err = tostring(l.err) | order by TimeGenerated desc | take 50'
    visualization: 'table'
  }
]

var workbookContent = {
  version: 'Notebook/1.0'
  items: concat(
    [
      {
        type: 1
        content: {
          json: '## ${appName}: operations overview\nRequest, latency, cold-start and error data from the app\'s structured JSON logs.'
        }
        name: 'header'
      }
    ],
    map(range(0, length(tiles)), i => {
      type: 3
      content: {
        version: 'KqlItem/1.0'
        query: tiles[i].query
        size: 0
        title: tiles[i].title
        timeContext: {
          durationMs: 604800000 // last 7 days
        }
        queryType: 0
        resourceType: 'microsoft.operationalinsights/workspaces'
        visualization: tiles[i].visualization
      }
      customWidth: '50'
      name: 'tile-${i}'
    })
  )
  fallbackResourceIds: [workspaceId]
  '$schema': 'https://github.com/Microsoft/Application-Insights-Workbooks/blob/master/schema/workbook.json'
}

resource workbook 'Microsoft.Insights/workbooks@2023-06-01' = {
  name: guid(resourceGroup().id, appName, 'ops-workbook')
  location: location
  tags: tags
  kind: 'shared'
  properties: {
    displayName: '${appName} operations'
    category: 'workbook'
    sourceId: workspaceId
    serializedData: string(workbookContent)
  }
}
