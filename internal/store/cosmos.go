package store

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

const (
	messagesContainer = "messages"
	statsContainer    = "stats"
	visitsID          = "visits"
)

type Cosmos struct {
	db       *azcosmos.DatabaseClient
	messages *azcosmos.ContainerClient
	stats    *azcosmos.ContainerClient
}

// NewCosmos authenticates with DefaultAzureCredential. In Container Apps that
// resolves to the user-assigned managed identity named by AZURE_CLIENT_ID; on a
// laptop it falls through to `az login`. No keys are ever handled.
func NewCosmos(endpoint, database string) (*Cosmos, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	client, err := azcosmos.NewClient(endpoint, cred, nil)
	if err != nil {
		return nil, err
	}

	db, err := client.NewDatabase(database)
	if err != nil {
		return nil, err
	}
	messages, err := db.NewContainer(messagesContainer)
	if err != nil {
		return nil, err
	}
	stats, err := db.NewContainer(statsContainer)
	if err != nil {
		return nil, err
	}

	return &Cosmos{db: db, messages: messages, stats: stats}, nil
}

func (c *Cosmos) SaveMessage(ctx context.Context, m Message) error {
	body, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = c.messages.CreateItem(ctx, azcosmos.NewPartitionKeyString(m.Kind), body, nil)
	return err
}

// IncrementVisits uses a server-side patch increment so concurrent replicas
// never lose an update. The document is created on first use.
func (c *Cosmos) IncrementVisits(ctx context.Context) (int64, error) {
	pk := azcosmos.NewPartitionKeyString(visitsID)
	ops := azcosmos.PatchOperations{}
	ops.AppendIncrement("/count", 1)
	opts := &azcosmos.ItemOptions{EnableContentResponseOnWrite: true}

	resp, err := c.stats.PatchItem(ctx, pk, visitsID, ops, opts)
	if isStatus(err, http.StatusNotFound) {
		doc, _ := json.Marshal(map[string]any{"id": visitsID, "count": 1})
		_, err = c.stats.CreateItem(ctx, pk, doc, nil)
		if !isStatus(err, http.StatusConflict) {
			return 1, err
		}
		// Another replica created it first; the increment will now succeed.
		resp, err = c.stats.PatchItem(ctx, pk, visitsID, ops, opts)
	}
	if err != nil {
		return 0, err
	}

	var out struct {
		Count int64 `json:"count"`
	}
	err = json.Unmarshal(resp.Value, &out)
	return out.Count, err
}

func (c *Cosmos) Ping(ctx context.Context) error {
	_, err := c.db.Read(ctx, nil)
	return err
}

func (c *Cosmos) Backend() string { return "Azure Cosmos DB (NoSQL, free tier)" }

func isStatus(err error, code int) bool {
	var respErr *azcore.ResponseError
	return errors.As(err, &respErr) && respErr.StatusCode == code
}
