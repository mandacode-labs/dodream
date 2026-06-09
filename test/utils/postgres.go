// Package testutil provides utilities for integration testing.
package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/lib/pq"

	"github.com/mandacode-labs/dodream/ent"
	"github.com/mandacode-labs/dodream/ent/enttest"
)

// NewPostgresContainer creates a new PostgreSQL container for testing.
func NewPostgresContainer(ctx context.Context, t *testing.T) (*postgres.PostgresContainer, string) {
	t.Helper()

	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("dodream_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForAll(
				wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
				wait.ForListeningPort("5432/tcp"),
			).WithDeadline(60*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	})

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	return container, connStr
}

// NewEntClient creates a new ent client with auto-migration for testing.
func NewEntClient(t *testing.T, connStr string) *ent.Client {
	t.Helper()

	return enttest.Open(t, "postgres", connStr)
}

// SetupTestDB is a convenience function that creates a postgres container and ent client.
func SetupTestDB(ctx context.Context, t *testing.T) (*ent.Client, func()) {
	t.Helper()

	_, connStr := NewPostgresContainer(ctx, t)
	client := NewEntClient(t, connStr)

	cleanup := func() {
		client.Close()
	}

	return client, cleanup
}
