package deck_test

import (
	"context"
	"testing"

	"github.com/mandacode-labs/dodream/internal/deck"
	"github.com/mandacode-labs/dodream/internal/testutil"
	"github.com/mandacode-labs/dodream/internal/user"
)

func TestStore_Create(t *testing.T) {
	ctx := context.Background()
	client, cleanup := testutil.SetupTestDB(ctx, t)
	defer cleanup()

	userStore := user.NewStore(client)
	u := user.New("testuser", "provider_123")
	createdUser, err := userStore.Create(ctx, u)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	store := deck.NewStore(client)
	d := deck.New("My Deck", createdUser.ID().String())

	created, err := store.Create(ctx, d)
	if err != nil {
		t.Fatalf("create deck: %v", err)
	}

	if created.Name() != "My Deck" {
		t.Errorf("name = %q, want %q", created.Name(), "My Deck")
	}
	if created.Creator() != createdUser.ID().String() {
		t.Errorf("creator = %q, want %q", created.Creator(), createdUser.ID().String())
	}
}

func TestStore_GetByID(t *testing.T) {
	ctx := context.Background()
	client, cleanup := testutil.SetupTestDB(ctx, t)
	defer cleanup()

	userStore := user.NewStore(client)
	u := user.New("testuser", "provider_123")
	createdUser, err := userStore.Create(ctx, u)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	store := deck.NewStore(client)
	d := deck.New("My Deck", createdUser.ID().String())
	created, err := store.Create(ctx, d)
	if err != nil {
		t.Fatalf("create deck: %v", err)
	}

	got, err := store.GetByID(ctx, created.ID())
	if err != nil {
		t.Fatalf("get deck by id: %v", err)
	}

	if got.Name() != created.Name() {
		t.Errorf("name = %q, want %q", got.Name(), created.Name())
	}
}

func TestStore_Update(t *testing.T) {
	ctx := context.Background()
	client, cleanup := testutil.SetupTestDB(ctx, t)
	defer cleanup()

	userStore := user.NewStore(client)
	u := user.New("testuser", "provider_123")
	createdUser, err := userStore.Create(ctx, u)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	store := deck.NewStore(client)
	d := deck.New("My Deck", createdUser.ID().String())
	created, err := store.Create(ctx, d)
	if err != nil {
		t.Fatalf("create deck: %v", err)
	}

	created.SetName("Updated Deck")

	updated, err := store.Update(ctx, created)
	if err != nil {
		t.Fatalf("update deck: %v", err)
	}

	if updated.Name() != "Updated Deck" {
		t.Errorf("name = %q, want %q", updated.Name(), "Updated Deck")
	}
}

func TestStore_Delete(t *testing.T) {
	ctx := context.Background()
	client, cleanup := testutil.SetupTestDB(ctx, t)
	defer cleanup()

	userStore := user.NewStore(client)
	u := user.New("testuser", "provider_123")
	createdUser, err := userStore.Create(ctx, u)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	store := deck.NewStore(client)
	d := deck.New("My Deck", createdUser.ID().String())
	created, err := store.Create(ctx, d)
	if err != nil {
		t.Fatalf("create deck: %v", err)
	}

	if err := store.Delete(ctx, created.ID()); err != nil {
		t.Fatalf("delete deck: %v", err)
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

	userStore := user.NewStore(client)
	u := user.New("testuser", "provider_123")
	createdUser, err := userStore.Create(ctx, u)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	store := deck.NewStore(client)
	for i := 0; i < 3; i++ {
		d := deck.New("Deck "+string(rune('0'+i)), createdUser.ID().String())
		if _, err := store.Create(ctx, d); err != nil {
			t.Fatalf("create deck: %v", err)
		}
	}

	decks, err := store.List(ctx)
	if err != nil {
		t.Fatalf("list decks: %v", err)
	}

	if len(decks) != 3 {
		t.Errorf("len(decks) = %d, want 3", len(decks))
	}
}

func TestStore_ListByCreator(t *testing.T) {
	ctx := context.Background()
	client, cleanup := testutil.SetupTestDB(ctx, t)
	defer cleanup()

	userStore := user.NewStore(client)
	u1 := user.New("user1", "provider_1")
	u2 := user.New("user2", "provider_2")
	createdUser1, err := userStore.Create(ctx, u1)
	if err != nil {
		t.Fatalf("create user1: %v", err)
	}
	createdUser2, err := userStore.Create(ctx, u2)
	if err != nil {
		t.Fatalf("create user2: %v", err)
	}

	store := deck.NewStore(client)
	for i := 0; i < 2; i++ {
		d := deck.New("Deck "+string(rune('0'+i)), createdUser1.ID().String())
		if _, err := store.Create(ctx, d); err != nil {
			t.Fatalf("create deck: %v", err)
		}
	}
	d := deck.New("Other Deck", createdUser2.ID().String())
	if _, err := store.Create(ctx, d); err != nil {
		t.Fatalf("create deck: %v", err)
	}

	decks, err := store.ListByCreator(ctx, createdUser1.ID().String())
	if err != nil {
		t.Fatalf("list decks by creator: %v", err)
	}

	if len(decks) != 2 {
		t.Errorf("len(decks) = %d, want 2", len(decks))
	}
}
