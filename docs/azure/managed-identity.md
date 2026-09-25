# Managed Identity and Entra ID RBAC

**What it is:** an identity in Microsoft Entra ID (formerly Azure AD) that Azure manages for a resource. Code running on that resource can get OAuth tokens for other Azure services, and **you never see or store a credential**. Azure rotates the underlying certificate automatically.

**Where:** [`infra/modules/identity.bicep`](../../infra/modules/identity.bicep), the `sqlRoleAssignments` in [`cosmos.bicep`](../../infra/modules/cosmos.bicep), and `NewCosmos` in [`internal/store/cosmos.go`](../../internal/store/cosmos.go)

## How it's wired

```
Container App ──attached──▶ user-assigned identity "id-say-hi"
                                   │ principalId
                                   ▼
            Cosmos DB sqlRoleAssignment: Built-in Data Contributor
            scope: /dbs/sayhi   (only this database, only data-plane)
```

1. Bicep creates the identity, then uses its `principalId` in a Cosmos DB role assignment, then attaches the identity to the container app.
2. The app's environment has `AZURE_CLIENT_ID=<identity client id>`.
3. `azidentity.NewDefaultAzureCredential()` in the Go SDK tries a chain of credential types. In Container Apps it finds the managed identity endpoint and uses the identity named by `AZURE_CLIENT_ID`. On a laptop it falls back to `az login`, so the **same code works everywhere** with no `if prod {...}` branches.

## Why user-assigned rather than system-assigned
A system-assigned identity only exists once its resource exists, and it's deleted along with that resource. Choosing user-assigned:
- lets the role assignment be created *before* the app's first revision starts, so the first replica never runs without database access;
- keeps permissions if the app is deleted and recreated.

## Least privilege
- The role is **data-plane only** (read and write items). The app can't change the account, keys, networking or throughput.
- The scope is **one database**, not the whole account.
- Account keys are disabled, so the identity can't be sidestepped with a leaked connection string.

## In production I would...
- Use separate identities per app and per environment, and review them with Entra access reviews.
- Add Azure Policy to deny Cosmos DB accounts or storage accounts where local (key) auth is enabled.
