package deck

import (
	"context"

	"github.com/mandacode-labs/dodream/internal/core"
	"github.com/mandacode-labs/dodream/internal/store"
)

// Service provides deck-related business logic.
type Service struct {
	store *store.DeckStore
}

// NewService creates a new Service with the given store.
func NewService(store *store.DeckStore) *Service {
	return &Service{store: store}
}

// Create creates a new deck.
func (s *Service) Create(ctx context.Context, name string, creator core.UserID) (*core.Deck, error) {
	deck := core.NewDeck(name, creator)
	return s.store.Create(ctx, deck)
}

// GetByID retrieves a deck by ID.
func (s *Service) GetByID(ctx context.Context, id core.DeckID) (*core.Deck, error) {
	return s.store.GetByID(ctx, id)
}

// List retrieves all decks.
func (s *Service) List(ctx context.Context) ([]*core.Deck, error) {
	return s.store.List(ctx)
}

// ListByCreator retrieves all decks created by a specific user.
func (s *Service) ListByCreator(ctx context.Context, userID core.UserID) ([]*core.Deck, error) {
	return s.store.ListByCreator(ctx, userID)
}
