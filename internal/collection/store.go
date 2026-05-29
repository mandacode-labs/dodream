package collection

import (
	"context"

	"github.com/mandacode-labs/dodream/ent"
	entcollection "github.com/mandacode-labs/dodream/ent/collection"
	entcollectioncard "github.com/mandacode-labs/dodream/ent/collectioncard"
	"github.com/mandacode-labs/dodream/ent/user"
	"github.com/mandacode-labs/dodream/internal/errs"
)

// Store provides database operations for collections.
type Store struct {
	client *ent.Client
}

// NewStore creates a new Store with the given ent client.
func NewStore(client *ent.Client) *Store {
	return &Store{client: client}
}

// Compile-time check that Store implements Repository.
var _ Repository = (*Store)(nil)

func mapEntError(op string, err error) error {
	switch {
	case ent.IsNotFound(err):
		return errs.Wrap(errs.ErrNotFound, op, err)
	case ent.IsConstraintError(err):
		return errs.Wrap(errs.ErrConflict, op, err)
	default:
		return errs.Wrap(errs.ErrInternal, op, err)
	}
}

// Create inserts a new collection into the database.
func (s *Store) Create(ctx context.Context, c *Collection) (*Collection, error) {
	const op = "create collection"
	created, err := s.client.Collection.Create().
		SetID(c.ID().String()).
		SetName(c.Name()).
		SetCreatorID(c.CreatorID()).
		Save(ctx)
	if err != nil {
		return nil, mapEntError(op, err)
	}
	return NewWithID(
		ID(created.ID),
		created.Name,
		c.CreatorID(),
		created.CreatedAt,
		created.UpdatedAt,
	), nil
}

// GetByID retrieves a collection by ID.
func (s *Store) GetByID(ctx context.Context, id ID) (*Collection, error) {
	const op = "get collection by id"
	c, err := s.client.Collection.Query().
		Where(entcollection.ID(id.String())).
		WithCreator().
		Only(ctx)
	if err != nil {
		return nil, mapEntError(op, err)
	}
	return NewWithID(
		ID(c.ID),
		c.Name,
		c.Edges.Creator.ID,
		c.CreatedAt,
		c.UpdatedAt,
	), nil
}

// Update modifies an existing collection.
func (s *Store) Update(ctx context.Context, c *Collection) (*Collection, error) {
	const op = "update collection"
	updated, err := s.client.Collection.UpdateOneID(c.ID().String()).
		SetName(c.Name()).
		Save(ctx)
	if err != nil {
		return nil, mapEntError(op, err)
	}
	return NewWithID(
		ID(updated.ID),
		updated.Name,
		c.CreatorID(),
		updated.CreatedAt,
		updated.UpdatedAt,
	), nil
}

// Delete removes a collection by ID.
func (s *Store) Delete(ctx context.Context, id ID) error {
	const op = "delete collection"
	err := s.client.Collection.DeleteOneID(id.String()).Exec(ctx)
	if err != nil {
		return mapEntError(op, err)
	}
	return nil
}

// List retrieves all collections.
func (s *Store) List(ctx context.Context) ([]*Collection, error) {
	const op = "list collections"
	collections, err := s.client.Collection.Query().WithCreator().All(ctx)
	if err != nil {
		return nil, mapEntError(op, err)
	}

	result := make([]*Collection, len(collections))
	for i, c := range collections {
		result[i] = NewWithID(
			ID(c.ID),
			c.Name,
			c.Edges.Creator.ID,
			c.CreatedAt,
			c.UpdatedAt,
		)
	}
	return result, nil
}

// ListByCreator retrieves all collections created by a specific user.
func (s *Store) ListByCreator(ctx context.Context, userID string) ([]*Collection, error) {
	const op = "list collections by creator"
	collections, err := s.client.Collection.Query().
		Where(entcollection.HasCreatorWith(user.ID(userID))).
		WithCreator().
		All(ctx)
	if err != nil {
		return nil, mapEntError(op, err)
	}

	result := make([]*Collection, len(collections))
	for i, c := range collections {
		result[i] = NewWithID(
			ID(c.ID),
			c.Name,
			c.Edges.Creator.ID,
			c.CreatedAt,
			c.UpdatedAt,
		)
	}
	return result, nil
}

// AddCard adds a card to a collection with an optional deck reference.
func (s *Store) AddCard(ctx context.Context, cc *CollectionCard) (*CollectionCard, error) {
	const op = "add card to collection"
	builder := s.client.CollectionCard.Create().
		SetID(cc.ID().String()).
		SetCollectionID(cc.CollectionID().String()).
		SetCardID(cc.CardID())

	if cc.DeckID() != nil {
		builder.SetDeckID(*cc.DeckID())
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return nil, mapEntError(op, err)
	}

	// Fetch with edges to return consistent data
	fetched, err := s.client.CollectionCard.Query().
		Where(entcollectioncard.ID(created.ID)).
		WithCollection().
		WithCard().
		WithDeck().
		Only(ctx)
	if err != nil {
		return nil, mapEntError("fetch collection card", err)
	}

	var resultDeckID *string
	if fetched.Edges.Deck != nil {
		id := fetched.Edges.Deck.ID
		resultDeckID = &id
	}

	return NewCollectionCardWithID(
		CollectionCardID(fetched.ID),
		ID(fetched.Edges.Collection.ID),
		fetched.Edges.Card.ID,
		resultDeckID,
		fetched.CreatedAt,
		fetched.UpdatedAt,
	), nil
}

// RemoveCard removes a card from a collection by the collection card ID.
func (s *Store) RemoveCard(ctx context.Context, collectionCardID CollectionCardID) error {
	const op = "remove card from collection"
	err := s.client.CollectionCard.DeleteOneID(collectionCardID.String()).Exec(ctx)
	if err != nil {
		return mapEntError(op, err)
	}
	return nil
}

// ListCards retrieves all cards within a specific collection.
func (s *Store) ListCards(ctx context.Context, collectionID ID) ([]*CollectionCard, error) {
	const op = "list collection cards"
	cards, err := s.client.CollectionCard.Query().
		Where(entcollectioncard.HasCollectionWith(entcollection.ID(collectionID.String()))).
		WithCollection().
		WithCard().
		WithDeck().
		All(ctx)
	if err != nil {
		return nil, mapEntError(op, err)
	}

	result := make([]*CollectionCard, len(cards))
	for i, c := range cards {
		var deckID *string
		if c.Edges.Deck != nil {
			id := c.Edges.Deck.ID
			deckID = &id
		}
		result[i] = NewCollectionCardWithID(
			CollectionCardID(c.ID),
			ID(c.Edges.Collection.ID),
			c.Edges.Card.ID,
			deckID,
			c.CreatedAt,
			c.UpdatedAt,
		)
	}
	return result, nil
}
