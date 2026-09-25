package store

import (
	"context"
	"sync"
)

type Memory struct {
	mu       sync.Mutex
	visits   int64
	Messages []Message
}

func NewMemory() *Memory {
	return &Memory{}
}

func (m *Memory) SaveMessage(_ context.Context, msg Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Messages = append(m.Messages, msg)
	return nil
}

func (m *Memory) IncrementVisits(_ context.Context) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.visits++
	return m.visits, nil
}

func (m *Memory) Ping(context.Context) error { return nil }

func (m *Memory) Backend() string { return "in-memory (local dev)" }
