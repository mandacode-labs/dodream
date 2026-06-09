//go:build integration

package integration

import (
	"testing"

	"github.com/mandacode-labs/dodream/internal/user"
)

func TestUserStore_Create(t *testing.T) {
	ctx, client, cleanup := SetupTestDB(t)
	defer cleanup()

	store := user.NewStore(client)
	u := user.New("testuser", "provider_123")

	created, err := store.Create(ctx, u)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	if created.Nickname() != "testuser" {
		t.Errorf("nickname = %q, want %q", created.Nickname(), "testuser")
	}
	if created.ProviderID() != "provider_123" {
		t.Errorf("providerID = %q, want %q", created.ProviderID(), "provider_123")
	}
}

func TestUserStore_GetByID(t *testing.T) {
	ctx, client, cleanup := SetupTestDB(t)
	defer cleanup()

	store := user.NewStore(client)
	u := user.New("testuser", "provider_123")

	created, err := store.Create(ctx, u)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	got, err := store.GetByID(ctx, created.ID())
	if err != nil {
		t.Fatalf("get user by id: %v", err)
	}

	if got.Nickname() != created.Nickname() {
		t.Errorf("nickname = %q, want %q", got.Nickname(), created.Nickname())
	}
}

func TestUserStore_Update(t *testing.T) {
	ctx, client, cleanup := SetupTestDB(t)
	defer cleanup()

	store := user.NewStore(client)
	u := user.New("testuser", "provider_123")

	created, err := store.Create(ctx, u)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	created.SetNickname("updated")
	created.SetProviderID("updated_provider")

	updated, err := store.Update(ctx, created)
	if err != nil {
		t.Fatalf("update user: %v", err)
	}

	if updated.Nickname() != "updated" {
		t.Errorf("nickname = %q, want %q", updated.Nickname(), "updated")
	}
	if updated.ProviderID() != "updated_provider" {
		t.Errorf("providerID = %q, want %q", updated.ProviderID(), "updated_provider")
	}
}

func TestUserStore_Delete(t *testing.T) {
	ctx, client, cleanup := SetupTestDB(t)
	defer cleanup()

	store := user.NewStore(client)
	u := user.New("testuser", "provider_123")

	created, err := store.Create(ctx, u)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	if err := store.Delete(ctx, created.ID()); err != nil {
		t.Fatalf("delete user: %v", err)
	}

	_, err = store.GetByID(ctx, created.ID())
	if err == nil {
		t.Fatal("expected error after deletion, got nil")
	}
}
