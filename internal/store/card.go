package store

import (
	"context"
	"fmt"

	"github.com/mandacode-labs/dodream/ent"
	"github.com/mandacode-labs/dodream/ent/card"
	"github.com/mandacode-labs/dodream/ent/deck"
	"github.com/mandacode-labs/dodream/internal/core"
)

// CardStore provides database operations for cards.
type CardStore struct {
	client *ent.Client
}

// NewCardStore creates a new CardStore with the given ent client.
func NewCardStore(client *ent.Client) *CardStore {
	return &CardStore{client: client}
}

// Create inserts a new card into the database.
func (s *CardStore) Create(ctx context.Context, c *core.Card) (*core.Card, error) {
	created, err := s.client.Card.Create().
		SetID(c.ID().String()).
		SetHint(c.Hint()).
		SetContent(c.Content()).
		SetCreatorID(c.Creator().String()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create card: %w", err)
	}
	return core.NewCard(
		created.Hint,
		created.Content,
		c.Creator(),
	), nil
}

// GetByID retrieves a card by their ID.
func (s *CardStore) GetByID(ctx context.Context, id core.CardID) (*core.Card, error) {
	c, err := s.client.Card.Query().
		Where(card.ID(id.String())).
		WithCreator().
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("get card by id: %w", err)
	}
	return core.NewCard(
		c.Hint,
		c.Content,
		core.UserID(c.Edges.Creator.ID),
	), nil
}

// Update modifies an existing card.
func (s *CardStore) Update(ctx context.Context, c *core.Card) (*core.Card, error) {
	updated, err := s.client.Card.UpdateOneID(c.ID().String()).
		SetHint(c.Hint()).
		SetContent(c.Content()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("update card: %w", err)
	}
	return core.NewCard(
		updated.Hint,
		updated.Content,
		c.Creator(),
	), nil
}

// Delete removes a card by their ID.
func (s *CardStore) Delete(ctx context.Context, id core.CardID) error {
	err := s.client.Card.DeleteOneID(id.String()).Exec(ctx)
	if err != nil {
		return fmt.Errorf("delete card: %w", err)
	}
	return nil
}

// List retrieves all cards.
func (s *CardStore) List(ctx context.Context) ([]*core.Card, error) {
	cards, err := s.client.Card.Query().WithCreator().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("list cards: %w", err)
	}

	result := make([]*core.Card, len(cards))
	for i, c := range cards {
		result[i] = core.NewCard(
			c.Hint,
			c.Content,
			core.UserID(c.Edges.Creator.ID),
		)
	}
	return result, nil
}

// AddToDeck adds a card to a deck.
func (s *CardStore) AddToDeck(ctx context.Context, cardID core.CardID, deckID core.DeckID) error {
	err := s.client.Deck.UpdateOneID(deckID.String()).
		AddCardIDs(cardID.String()).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("add card to deck: %w", err)
	}
	return nil
}

// RemoveFromDeck removes a card from a deck.
func (s *CardStore) RemoveFromDeck(ctx context.Context, cardID core.CardID, deckID core.DeckID) error {
	err := s.client.Deck.UpdateOneID(deckID.String()).
		RemoveCardIDs(cardID.String()).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("remove card from deck: %w", err)
	}
	return nil
}

// ListByDeck retrieves all cards in a specific deck.
func (s *CardStore) ListByDeck(ctx context.Context, deckID core.DeckID) ([]*core.Card, error) {
	cards, err := s.client.Card.Query().
		Where(card.HasDecksWith(deck.ID(deckID.String()))).
		WithCreator().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("list cards by deck: %w", err)
	}

	result := make([]*core.Card, len(cards))
	for i, c := range cards {
		result[i] = core.NewCard(
			c.Hint,
			c.Content,
			core.UserID(c.Edges.Creator.ID),
		)
	}
	return result, nil
}
