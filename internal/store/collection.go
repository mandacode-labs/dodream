package store

import (
	"context"
	"fmt"

	"github.com/mandacode-labs/dodream/ent"
	"github.com/mandacode-labs/dodream/ent/collection"
	"github.com/mandacode-labs/dodream/ent/collectioncard"
	"github.com/mandacode-labs/dodream/ent/user"
	"github.com/mandacode-labs/dodream/internal/core"
)

// CollectionStore provides database operations for collections.
type CollectionStore struct {
	client *ent.Client
}

// NewCollectionStore creates a new CollectionStore with the given ent client.
func NewCollectionStore(client *ent.Client) *CollectionStore {
	return &CollectionStore{client: client}
}

// Create inserts a new collection into the database.
func (s *CollectionStore) Create(ctx context.Context, c *core.Collection) (*core.Collection, error) {
	created, err := s.client.Collection.Create().
		SetID(c.ID().String()).
		SetName(c.Name()).
		SetCreatorID(c.Creator().String()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create collection: %w", err)
	}
	return core.NewCollection(
		core.CollectionID(created.ID),
		created.Name,
		c.Creator(),
	), nil
}

// GetByID retrieves a collection by their ID.
func (s *CollectionStore) GetByID(ctx context.Context, id core.CollectionID) (*core.Collection, error) {
	c, err := s.client.Collection.Query().
		Where(collection.ID(id.String())).
		WithCreator().
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("get collection by id: %w", err)
	}
	return core.NewCollection(
		core.CollectionID(c.ID),
		c.Name,
		core.UserID(c.Edges.Creator.ID),
	), nil
}

// Update modifies an existing collection.
func (s *CollectionStore) Update(ctx context.Context, c *core.Collection) (*core.Collection, error) {
	updated, err := s.client.Collection.UpdateOneID(c.ID().String()).
		SetName(c.Name()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("update collection: %w", err)
	}
	return core.NewCollection(
		core.CollectionID(updated.ID),
		updated.Name,
		c.Creator(),
	), nil
}

// Delete removes a collection by their ID.
func (s *CollectionStore) Delete(ctx context.Context, id core.CollectionID) error {
	err := s.client.Collection.DeleteOneID(id.String()).Exec(ctx)
	if err != nil {
		return fmt.Errorf("delete collection: %w", err)
	}
	return nil
}

// List retrieves all collections.
func (s *CollectionStore) List(ctx context.Context) ([]*core.Collection, error) {
	collections, err := s.client.Collection.Query().WithCreator().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("list collections: %w", err)
	}

	result := make([]*core.Collection, len(collections))
	for i, c := range collections {
		result[i] = core.NewCollection(
			core.CollectionID(c.ID),
			c.Name,
			core.UserID(c.Edges.Creator.ID),
		)
	}
	return result, nil
}

// ListByCreator retrieves all collections created by a specific user.
func (s *CollectionStore) ListByCreator(ctx context.Context, userID core.UserID) ([]*core.Collection, error) {
	collections, err := s.client.Collection.Query().
		Where(collection.HasCreatorWith(user.ID(userID.String()))).
		WithCreator().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("list collections by creator: %w", err)
	}

	result := make([]*core.Collection, len(collections))
	for i, c := range collections {
		result[i] = core.NewCollection(
			core.CollectionID(c.ID),
			c.Name,
			core.UserID(c.Edges.Creator.ID),
		)
	}
	return result, nil
}

// AddCard adds a card to a collection with an optional deck reference.
func (s *CollectionStore) AddCard(ctx context.Context, collectionID core.CollectionID, cardID core.CardID, deckID *core.DeckID) (*core.CollectionCard, error) {
	builder := s.client.CollectionCard.Create().
		SetCollectionID(collectionID.String()).
		SetCardID(cardID.String())

	if deckID != nil {
		builder.SetDeckID(deckID.String())
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("add card to collection: %w", err)
	}

	var dID *core.DeckID
	if created.Edges.Deck != nil {
		id := core.DeckID(created.Edges.Deck.ID)
		dID = &id
	}

	return core.NewCollectionCard(
		core.CollectionCardID(created.ID),
		core.CollectionID(created.Edges.Collection.ID),
		core.CardID(created.Edges.Card.ID),
		dID,
	), nil
}

// RemoveCard removes a card from a collection by the collection card ID.
func (s *CollectionStore) RemoveCard(ctx context.Context, collectionCardID core.CollectionCardID) error {
	err := s.client.CollectionCard.DeleteOneID(collectionCardID.String()).Exec(ctx)
	if err != nil {
		return fmt.Errorf("remove card from collection: %w", err)
	}
	return nil
}

// ListCards retrieves all cards within a specific collection.
func (s *CollectionStore) ListCards(ctx context.Context, collectionID core.CollectionID) ([]*core.CollectionCard, error) {
	cards, err := s.client.CollectionCard.Query().
		Where(collectioncard.HasCollectionWith(collection.ID(collectionID.String()))).
		WithCard().
		WithDeck().
		WithCollection().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("list collection cards: %w", err)
	}

	result := make([]*core.CollectionCard, len(cards))
	for i, c := range cards {
		var dID *core.DeckID
		if c.Edges.Deck != nil {
			id := core.DeckID(c.Edges.Deck.ID)
			dID = &id
		}
		result[i] = core.NewCollectionCard(
			core.CollectionCardID(c.ID),
			core.CollectionID(c.Edges.Collection.ID),
			core.CardID(c.Edges.Card.ID),
			dID,
		)
	}
	return result, nil
}
