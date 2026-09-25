// Package runtimeinfo reports where and how this replica is running. Azure
// Container Apps injects CONTAINER_APP_* variables into every replica; the
// region and version are passed in by the Bicep template and the CI pipeline.
package runtimeinfo

import (
	"os"
	"sync/atomic"
	"time"
)

var (
	started  = time.Now()
	requests atomic.Int64
)

type Info struct {
	Platform string
	App      string
	Revision string
	Replica  string
	Region   string
	Version  string
	Uptime   string
	Started  string
	Requests int64
}

// CountRequest records a request served by this replica.
func CountRequest() {
	requests.Add(1)
}

func Current() Info {
	platform := "local"
	if os.Getenv("CONTAINER_APP_NAME") != "" {
		platform = "Azure Container Apps"
	}

	return Info{
		Platform: platform,
		App:      envOr("CONTAINER_APP_NAME", "say-hi"),
		Revision: envOr("CONTAINER_APP_REVISION", "n/a"),
		Replica:  envOr("CONTAINER_APP_REPLICA_NAME", hostname()),
		Region:   envOr("AZURE_REGION", "n/a"),
		Version:  envOr("APP_VERSION", "dev"),
		Uptime:   time.Since(started).Round(time.Second).String(),
		Started:  started.UTC().Format(time.RFC3339),
		Requests: requests.Load(),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func hostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return h
}
