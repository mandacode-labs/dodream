package store

import (
	"context"
	"fmt"

	"github.com/mandacode-labs/dodream/ent"
	"github.com/mandacode-labs/dodream/ent/deck"
	"github.com/mandacode-labs/dodream/ent/user"
	"github.com/mandacode-labs/dodream/internal/core"
)

// DeckStore provides database operations for decks.
type DeckStore struct {
	client *ent.Client
}

// NewDeckStore creates a new DeckStore with the given ent client.
func NewDeckStore(client *ent.Client) *DeckStore {
	return &DeckStore{client: client}
}

// Create inserts a new deck into the database.
func (s *DeckStore) Create(ctx context.Context, d *core.Deck) (*core.Deck, error) {
	created, err := s.client.Deck.Create().
		SetID(d.ID().String()).
		SetName(d.Name()).
		SetCreatorID(d.Creator().String()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create deck: %w", err)
	}
	return core.NewDeck(
		core.DeckID(created.ID),
		created.Name,
		d.Creator(),
	), nil
}

// GetByID retrieves a deck by their ID.
func (s *DeckStore) GetByID(ctx context.Context, id core.DeckID) (*core.Deck, error) {
	d, err := s.client.Deck.Query().
		Where(deck.ID(id.String())).
		WithCreator().
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("get deck by id: %w", err)
	}
	return core.NewDeck(
		core.DeckID(d.ID),
		d.Name,
		core.UserID(d.Edges.Creator.ID),
	), nil
}

// Update modifies an existing deck.
func (s *DeckStore) Update(ctx context.Context, d *core.Deck) (*core.Deck, error) {
	updated, err := s.client.Deck.UpdateOneID(d.ID().String()).
		SetName(d.Name()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("update deck: %w", err)
	}
	return core.NewDeck(
		core.DeckID(updated.ID),
		updated.Name,
		d.Creator(),
	), nil
}

// Delete removes a deck by their ID.
func (s *DeckStore) Delete(ctx context.Context, id core.DeckID) error {
	err := s.client.Deck.DeleteOneID(id.String()).Exec(ctx)
	if err != nil {
		return fmt.Errorf("delete deck: %w", err)
	}
	return nil
}

// List retrieves all decks.
func (s *DeckStore) List(ctx context.Context) ([]*core.Deck, error) {
	decks, err := s.client.Deck.Query().WithCreator().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("list decks: %w", err)
	}

	result := make([]*core.Deck, len(decks))
	for i, d := range decks {
		result[i] = core.NewDeck(
			core.DeckID(d.ID),
			d.Name,
			core.UserID(d.Edges.Creator.ID),
		)
	}
	return result, nil
}

// ListByCreator retrieves all decks created by a specific user.
func (s *DeckStore) ListByCreator(ctx context.Context, userID core.UserID) ([]*core.Deck, error) {
	decks, err := s.client.Deck.Query().
		Where(deck.HasCreatorWith(user.ID(userID.String()))).
		WithCreator().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("list decks by creator: %w", err)
	}

	result := make([]*core.Deck, len(decks))
	for i, d := range decks {
		result[i] = core.NewDeck(
			core.DeckID(d.ID),
			d.Name,
			core.UserID(d.Edges.Creator.ID),
		)
	}
	return result, nil
}
