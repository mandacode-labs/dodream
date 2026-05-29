//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/mandacode-labs/dodream/ent"
	testutil "github.com/mandacode-labs/dodream/test/utils"
)

// SetupTestDB is a helper that sets up a test database for integration tests.
// It returns the ent client and a cleanup function.
func SetupTestDB(t *testing.T) (context.Context, *ent.Client, func()) {
	ctx := context.Background()
	client, cleanup := testutil.SetupTestDB(ctx, t)

	return ctx, client, func() {
		cleanup()
		client.Close()
	}
}
