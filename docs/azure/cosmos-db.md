# Azure Cosmos DB for NoSQL

**What it is:** Azure's globally distributed NoSQL database, with single-digit-millisecond latency and a choice of five consistency levels. Throughput is provisioned in **Request Units per second (RU/s)**, a normalised measure of CPU, IO and memory per operation.

**Where:** [`infra/modules/cosmos.bicep`](../../infra/modules/cosmos.bicep), [`internal/store/cosmos.go`](../../internal/store/cosmos.go)

## Data model

| Container | Partition key | Documents |
|---|---|---|
| `messages` | `/kind` | One per contact-form submission: `{id, kind:"contact", name, email, body, createdAt}` |
| `stats` | `/id` | A single `{id:"visits", count}` document |

The partition keys are chosen for a tiny dataset where every query targets a known partition. At scale, `messages` would be partitioned by something with high cardinality (for example, a date bucket) to avoid a hot partition.

## Things worth pointing out

- **Atomic counter without read-modify-write.** The visitor counter uses a *partial document update* (`PatchItem` with an `increment` operation). The increment runs on the server, so two replicas can't overwrite each other's update and no optimistic-concurrency retry loop is needed. The first increment creates the document, and a `409 Conflict` (another replica created it first) is handled.
- **Keys are disabled** (`disableLocalAuth: true`). The account can't be reached with a primary key or a connection string. The only way in is a Microsoft Entra ID token checked against **data-plane RBAC**:
  - the app's managed identity has *Built-in Data Contributor*, scoped to the `sayhi` database only;
  - the owner can optionally be granted *Built-in Data Reader* to browse messages in Data Explorer.
- **Throughput hard cap.** `capacity.totalThroughputLimit: 1000` makes Azure *reject* any change that would provision more than the free 1000 RU/s. Budgets only alert after the fact; this cap prevents the charge in the first place.
- **Shared database throughput.** The 1000 RU/s is provisioned on the database and shared by both containers, instead of each container needing its own minimum of 400 RU/s.
- **Session consistency.** Each client reads its own writes, and it's the cheapest option for RU cost.
- **Timeouts.** Every call has a context deadline (2–5 s). If Cosmos is slow, the visitor counter is hidden and the contact form tells the visitor to use email. Nothing else on the page is affected.

## Free tier
- Always free: the first **1000 RU/s and 25 GB** on one account per subscription.
- It must be enabled **when the account is created**. It can't be switched on later.
- It isn't compatible with serverless accounts. That's why this setup uses provisioned throughput.

## Local development
When `COSMOS_ENDPOINT` is unset, `store.FromEnv()` returns an in-memory store, so `go run ./cmd` works offline. To use the real database locally, set `COSMOS_ENDPOINT`, run `az login` and grant yourself the data-plane role. `DefaultAzureCredential` falls back to your CLI login.

## In production I would...
- Use private endpoints and `publicNetworkAccess: Disabled`.
- Enable **continuous backup** (point-in-time restore) and configure a TTL or retention policy for personal data such as contact emails (GDPR/CCPA).
- Choose partition keys from real access patterns, and alert on 429 (throttled) responses and normalised RU consumption.
