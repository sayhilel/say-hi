// Package store persists contact-form messages and the visitor counter.
//
// In Azure the backing store is Cosmos DB for NoSQL (free tier), reached with
// the container app's managed identity. Locally, when COSMOS_ENDPOINT is not
// set, an in-memory store is used so the site runs without any cloud access.
package store

import (
	"context"
	"os"
	"time"
)

type Message struct {
	ID        string    `json:"id"`
	Kind      string    `json:"kind"` // partition key; always "contact"
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
}

type Store interface {
	SaveMessage(ctx context.Context, m Message) error
	IncrementVisits(ctx context.Context) (int64, error)
	Ping(ctx context.Context) error
	Backend() string
}

// FromEnv returns a Cosmos-backed store when COSMOS_ENDPOINT is set and an
// in-memory store otherwise.
func FromEnv() (Store, error) {
	endpoint := os.Getenv("COSMOS_ENDPOINT")
	if endpoint == "" {
		return NewMemory(), nil
	}

	db := os.Getenv("COSMOS_DATABASE")
	if db == "" {
		db = "sayhi"
	}

	return NewCosmos(endpoint, db)
}
