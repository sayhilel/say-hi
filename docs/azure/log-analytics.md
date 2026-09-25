# Log Analytics, KQL and Workbooks

**What it is:** Log Analytics is the log store behind Azure Monitor. You query it with **KQL (Kusto Query Language)**. **Workbooks** are Azure Monitor's parameterised dashboards, built from KQL queries and stored as Azure resources.

**Where:** [`infra/modules/logs.bicep`](../../infra/modules/logs.bicep), [`infra/modules/workbook.bicep`](../../infra/modules/workbook.bicep), the `requestLogger` in [`cmd/main.go`](../../cmd/main.go)

## Structured logging
The app uses Go's standard `log/slog` with a JSON handler. Each request produces one line:

```json
{"time":"2026-09-24T20:16:38Z","level":"INFO","msg":"request","method":"POST",
 "path":"/command","status":200,"durationMs":1,"ip":"203.0.113.7","userAgent":"..."}
```

Container Apps collects stdout into the `ContainerAppConsoleLogs_CL` table, with the raw line in `Log_s`. Because the line is JSON, KQL can `parse_json` it and query every field. There's no log agent or SDK to maintain.

`/healthz` probe requests are deliberately *not* logged. Probes arrive every few seconds and would waste the free ingestion allowance.

## Useful queries
Open *Log Analytics workspace → Logs* and paste:

```kusto
// Requests per hour
ContainerAppConsoleLogs_CL
| where ContainerAppName_s == "say-hi"
| extend l = parse_json(Log_s)
| where l.msg == "request"
| summarize requests = count() by bin(TimeGenerated, 1h)
| render timechart
```

```kusto
// p50 / p95 latency by path
ContainerAppConsoleLogs_CL
| where ContainerAppName_s == "say-hi"
| extend l = parse_json(Log_s)
| where l.msg == "request"
| summarize p50 = percentile(toint(l.durationMs), 50),
            p95 = percentile(toint(l.durationMs), 95),
            hits = count() by path = tostring(l.path)
| order by hits desc
```

```kusto
// Errors with context
ContainerAppConsoleLogs_CL
| where ContainerAppName_s == "say-hi"
| extend l = parse_json(Log_s)
| where l.level == "ERROR"
| project TimeGenerated, RevisionName_s, msg = tostring(l.msg), err = tostring(l.err)
| order by TimeGenerated desc
```

```kusto
// Platform events: revision provisioning, image pulls, probe failures, scaling
ContainerAppSystemLogs_CL
| where ContainerAppName_s == "say-hi"
| project TimeGenerated, Reason_s, Log_s
| order by TimeGenerated desc
```

## The workbook
Bicep deploys a workbook called **"say-hi operations"**. You'll find it under *Monitor → Workbooks* or in the resource group. It has these tiles:
- requests per hour
- p50/p95 latency
- status-code mix
- top paths
- cold starts (the app logs `listening` once per replica start)
- contact messages received
- recent errors

The whole dashboard is stored as code in `workbook.bicep`, so it's versioned and reviewed like everything else.

## Cost controls
- The first **5 GB per month** of ingestion is free, and **31 days** of retention are included.
- `workspaceCapping.dailyQuotaGb: 0.15` stops ingestion for the rest of the day if it hits about 150 MB. About 30 days × 0.15 GB = 4.5 GB, which stays under the free allowance even in the worst case.

## In production I would...
- Add **Application Insights** with OpenTelemetry for distributed traces and dependency tracking. Its data lands in the same workspace.
- Add **alert rules**, for example 5xx rate > 1% or p95 > 1 s, routed to an action group (email, Teams, PagerDuty). Log search alerts are billed, so they're left out here.
- Use longer retention or archive tiers to meet regulatory record-keeping requirements.
