package card_test

import (
	"context"
	"testing"

	"github.com/mandacode-labs/dodream/internal/card"
	"github.com/mandacode-labs/dodream/internal/testutil"
	"github.com/mandacode-labs/dodream/internal/user"
)

func TestStore_Create(t *testing.T) {
	ctx := context.Background()
	client, cleanup := testutil.SetupTestDB(ctx, t)
	defer cleanup()

	// Create user first
	userStore := user.NewStore(client)
	u := user.New("testuser", "provider_123")
	createdUser, err := userStore.Create(ctx, u)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	store := card.NewStore(client)
	c := card.New("What is Go?", "A programming language", "Go is ...", createdUser.ID().String())

	created, err := store.Create(ctx, c)
	if err != nil {
		t.Fatalf("create card: %v", err)
	}

	if created.Question() != "What is Go?" {
		t.Errorf("question = %q, want %q", created.Question(), "What is Go?")
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

	store := card.NewStore(client)
	c := card.New("Q1", "H1", "C1", createdUser.ID().String())
	created, err := store.Create(ctx, c)
	if err != nil {
		t.Fatalf("create card: %v", err)
	}

	got, err := store.GetByID(ctx, created.ID())
	if err != nil {
		t.Fatalf("get card by id: %v", err)
	}

	if got.Question() != created.Question() {
		t.Errorf("question = %q, want %q", got.Question(), created.Question())
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

	store := card.NewStore(client)
	c := card.New("Q1", "H1", "C1", createdUser.ID().String())
	created, err := store.Create(ctx, c)
	if err != nil {
		t.Fatalf("create card: %v", err)
	}

	created.SetQuestion("Updated Q")
	created.SetHint("Updated H")
	created.SetContent("Updated C")

	updated, err := store.Update(ctx, created)
	if err != nil {
		t.Fatalf("update card: %v", err)
	}

	if updated.Question() != "Updated Q" {
		t.Errorf("question = %q, want %q", updated.Question(), "Updated Q")
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

	store := card.NewStore(client)
	c := card.New("Q1", "H1", "C1", createdUser.ID().String())
	created, err := store.Create(ctx, c)
	if err != nil {
		t.Fatalf("create card: %v", err)
	}

	if err := store.Delete(ctx, created.ID()); err != nil {
		t.Fatalf("delete card: %v", err)
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

	store := card.NewStore(client)
	for i := 0; i < 3; i++ {
		c := card.New("Q"+string(rune('0'+i)), "H", "C", createdUser.ID().String())
		if _, err := store.Create(ctx, c); err != nil {
			t.Fatalf("create card: %v", err)
		}
	}

	cards, err := store.List(ctx)
	if err != nil {
		t.Fatalf("list cards: %v", err)
	}

	if len(cards) != 3 {
		t.Errorf("len(cards) = %d, want 3", len(cards))
	}
}
