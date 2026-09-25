# Cost Management and FinOps guardrails

**What it is:** Azure Cost Management shows spend by resource, tag and service, and lets you set **budgets** that alert by email or trigger automation at thresholds.

**Where:** [`infra/modules/budget.bicep`](../../infra/modules/budget.bicep), plus the caps in the other modules

## Layered controls
The site has to cost $0 after the free-account credit expires, so cost is controlled at several layers. Prevention comes first; alerting is the backstop.

| Layer | Control | Type |
|---|---|---|
| Compute | `minReplicas: 0`, `maxReplicas: 1`, 0.25 vCPU | Prevents spend |
| Database | Free tier + `totalThroughputLimit: 1000` RU/s | Prevents spend (Azure rejects anything more) |
| Logging | Daily ingestion cap of 0.15 GB, probe logs suppressed | Prevents spend |
| Registry | GHCR instead of ACR | Avoids a fixed ~$5/month charge |
| Everything | $1 monthly budget: email at 50% and 100% actual and 100% forecast | Detects spend |
| Everything | Tags `app`, `managedBy`, `costCenter` | Makes cost reports attributable |

The budget is created when the `BUDGET_EMAIL` repository variable is set.

## Checking it
*Cost Management → Cost analysis*, scoped to `rg-say-hi` and grouped by *Service name*, should show $0.00. *Cost Management → Budgets* shows `budget-say-hi` and its thresholds.

## Tearing it all down
Everything lives in one resource group:
```bash
az group delete --name rg-say-hi
```

## In production I would...
- Use budgets with an **action group** that calls automation, for example scaling non-production to zero or notifying a Teams channel.
- Use Azure Advisor cost recommendations, reservations or savings plans for steady workloads, and showback/chargeback by `costCenter` tag.
