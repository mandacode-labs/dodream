package collection_test

import (
	"context"
	"testing"

	"github.com/mandacode-labs/dodream/internal/card"
	"github.com/mandacode-labs/dodream/internal/collection"
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

	store := collection.NewStore(client)
	c := collection.New("My Collection", createdUser.ID().String())

	created, err := store.Create(ctx, c)
	if err != nil {
		t.Fatalf("create collection: %v", err)
	}

	if created.Name() != "My Collection" {
		t.Errorf("name = %q, want %q", created.Name(), "My Collection")
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

	store := collection.NewStore(client)
	c := collection.New("My Collection", createdUser.ID().String())
	created, err := store.Create(ctx, c)
	if err != nil {
		t.Fatalf("create collection: %v", err)
	}

	got, err := store.GetByID(ctx, created.ID())
	if err != nil {
		t.Fatalf("get collection by id: %v", err)
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

	store := collection.NewStore(client)
	c := collection.New("My Collection", createdUser.ID().String())
	created, err := store.Create(ctx, c)
	if err != nil {
		t.Fatalf("create collection: %v", err)
	}

	created.SetName("Updated Collection")

	updated, err := store.Update(ctx, created)
	if err != nil {
		t.Fatalf("update collection: %v", err)
	}

	if updated.Name() != "Updated Collection" {
		t.Errorf("name = %q, want %q", updated.Name(), "Updated Collection")
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

	store := collection.NewStore(client)
	c := collection.New("My Collection", createdUser.ID().String())
	created, err := store.Create(ctx, c)
	if err != nil {
		t.Fatalf("create collection: %v", err)
	}

	if err := store.Delete(ctx, created.ID()); err != nil {
		t.Fatalf("delete collection: %v", err)
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

	store := collection.NewStore(client)
	for i := 0; i < 3; i++ {
		c := collection.New("Collection "+string(rune('0'+i)), createdUser.ID().String())
		if _, err := store.Create(ctx, c); err != nil {
			t.Fatalf("create collection: %v", err)
		}
	}

	collections, err := store.List(ctx)
	if err != nil {
		t.Fatalf("list collections: %v", err)
	}

	if len(collections) != 3 {
		t.Errorf("len(collections) = %d, want 3", len(collections))
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

	store := collection.NewStore(client)
	for i := 0; i < 2; i++ {
		c := collection.New("Collection "+string(rune('0'+i)), createdUser1.ID().String())
		if _, err := store.Create(ctx, c); err != nil {
			t.Fatalf("create collection: %v", err)
		}
	}
	c := collection.New("Other Collection", createdUser2.ID().String())
	if _, err := store.Create(ctx, c); err != nil {
		t.Fatalf("create collection: %v", err)
	}

	collections, err := store.ListByCreator(ctx, createdUser1.ID().String())
	if err != nil {
		t.Fatalf("list collections by creator: %v", err)
	}

	if len(collections) != 2 {
		t.Errorf("len(collections) = %d, want 2", len(collections))
	}
}

func TestStore_AddCard(t *testing.T) {
	ctx := context.Background()
	client, cleanup := testutil.SetupTestDB(ctx, t)
	defer cleanup()

	userStore := user.NewStore(client)
	u := user.New("testuser", "provider_123")
	createdUser, err := userStore.Create(ctx, u)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	cardStore := card.NewStore(client)
	c := card.New("Q1", "H1", "C1", createdUser.ID().String())
	createdCard, err := cardStore.Create(ctx, c)
	if err != nil {
		t.Fatalf("create card: %v", err)
	}

	store := collection.NewStore(client)
	col := collection.New("My Collection", createdUser.ID().String())
	createdCol, err := store.Create(ctx, col)
	if err != nil {
		t.Fatalf("create collection: %v", err)
	}

	cc := collection.NewCollectionCard(createdCol.ID(), createdCard.ID().String(), nil)
	createdCC, err := store.AddCard(ctx, cc)
	if err != nil {
		t.Fatalf("add card to collection: %v", err)
	}

	if createdCC.CardID() != createdCard.ID().String() {
		t.Errorf("cardID = %q, want %q", createdCC.CardID(), createdCard.ID().String())
	}
}

func TestStore_ListCards(t *testing.T) {
	ctx := context.Background()
	client, cleanup := testutil.SetupTestDB(ctx, t)
	defer cleanup()

	userStore := user.NewStore(client)
	u := user.New("testuser", "provider_123")
	createdUser, err := userStore.Create(ctx, u)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	cardStore := card.NewStore(client)
	c := card.New("Q1", "H1", "C1", createdUser.ID().String())
	createdCard, err := cardStore.Create(ctx, c)
	if err != nil {
		t.Fatalf("create card: %v", err)
	}

	store := collection.NewStore(client)
	col := collection.New("My Collection", createdUser.ID().String())
	createdCol, err := store.Create(ctx, col)
	if err != nil {
		t.Fatalf("create collection: %v", err)
	}

	cc := collection.NewCollectionCard(createdCol.ID(), createdCard.ID().String(), nil)
	if _, err := store.AddCard(ctx, cc); err != nil {
		t.Fatalf("add card to collection: %v", err)
	}

	cards, err := store.ListCards(ctx, createdCol.ID())
	if err != nil {
		t.Fatalf("list cards: %v", err)
	}

	if len(cards) != 1 {
		t.Errorf("len(cards) = %d, want 1", len(cards))
	}
}

func TestStore_AddCard_WithDeck(t *testing.T) {
	ctx := context.Background()
	client, cleanup := testutil.SetupTestDB(ctx, t)
	defer cleanup()

	userStore := user.NewStore(client)
	u := user.New("testuser", "provider_123")
	createdUser, err := userStore.Create(ctx, u)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	cardStore := card.NewStore(client)
	c := card.New("Q1", "H1", "C1", createdUser.ID().String())
	createdCard, err := cardStore.Create(ctx, c)
	if err != nil {
		t.Fatalf("create card: %v", err)
	}

	deckStore := deck.NewStore(client)
	d := deck.New("My Deck", createdUser.ID().String())
	createdDeck, err := deckStore.Create(ctx, d)
	if err != nil {
		t.Fatalf("create deck: %v", err)
	}

	store := collection.NewStore(client)
	col := collection.New("My Collection", createdUser.ID().String())
	createdCol, err := store.Create(ctx, col)
	if err != nil {
		t.Fatalf("create collection: %v", err)
	}

	deckID := createdDeck.ID().String()
	cc := collection.NewCollectionCard(createdCol.ID(), createdCard.ID().String(), &deckID)
	createdCC, err := store.AddCard(ctx, cc)
	if err != nil {
		t.Fatalf("add card to collection: %v", err)
	}

	if createdCC.DeckID() == nil || *createdCC.DeckID() != deckID {
		t.Errorf("deckID = %v, want %q", createdCC.DeckID(), deckID)
	}
}
