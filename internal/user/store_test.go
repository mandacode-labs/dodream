package user_test

import (
	"context"
	"testing"

	"github.com/mandacode-labs/dodream/internal/testutil"
	"github.com/mandacode-labs/dodream/internal/user"
)

func TestStore_Create(t *testing.T) {
	ctx := context.Background()
	client, cleanup := testutil.SetupTestDB(ctx, t)
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

func TestStore_GetByID(t *testing.T) {
	ctx := context.Background()
	client, cleanup := testutil.SetupTestDB(ctx, t)
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

func TestStore_Update(t *testing.T) {
	ctx := context.Background()
	client, cleanup := testutil.SetupTestDB(ctx, t)
	defer cleanup()

	store := user.NewStore(client)
	u := user.New("testuser", "provider_123")

	created, err := store.Create(ctx, u)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	created.SetNickname("updateduser")
	created.SetProviderID("provider_456")

	updated, err := store.Update(ctx, created)
	if err != nil {
		t.Fatalf("update user: %v", err)
	}

	if updated.Nickname() != "updateduser" {
		t.Errorf("nickname = %q, want %q", updated.Nickname(), "updateduser")
	}
	if updated.ProviderID() != "provider_456" {
		t.Errorf("providerID = %q, want %q", updated.ProviderID(), "provider_456")
	}
}

func TestStore_Delete(t *testing.T) {
	ctx := context.Background()
	client, cleanup := testutil.SetupTestDB(ctx, t)
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
		t.Error("expected error after delete, got nil")
	}
}

func TestStore_List(t *testing.T) {
	ctx := context.Background()
	client, cleanup := testutil.SetupTestDB(ctx, t)
	defer cleanup()

	store := user.NewStore(client)

	for i := 0; i < 3; i++ {
		u := user.New("user_"+string(rune('a'+i)), "provider_"+string(rune('0'+i)))
		if _, err := store.Create(ctx, u); err != nil {
			t.Fatalf("create user: %v", err)
		}
	}

	users, err := store.List(ctx)
	if err != nil {
		t.Fatalf("list users: %v", err)
	}

	if len(users) != 3 {
		t.Errorf("len(users) = %d, want 3", len(users))
	}
}
