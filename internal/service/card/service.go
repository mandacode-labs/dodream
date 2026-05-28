package card

import (
	"context"

	"github.com/mandacode-labs/dodream/internal/core"
	"github.com/mandacode-labs/dodream/internal/store"
)

// Service provides card-related business logic.
type Service struct {
	store *store.CardStore
}

// NewService creates a new Service with the given store.
func NewService(store *store.CardStore) *Service {
	return &Service{store: store}
}

// Create creates a new card.
func (s *Service) Create(ctx context.Context, hint string, content string, creator core.UserID) (*core.Card, error) {
	card := core.NewCard(hint, content, creator)
	return s.store.Create(ctx, card)
}

// GetByID retrieves a card by ID.
func (s *Service) GetByID(ctx context.Context, id core.CardID) (*core.Card, error) {
	return s.store.GetByID(ctx, id)
}

// List retrieves all cards.
func (s *Service) List(ctx context.Context) ([]*core.Card, error) {
	return s.store.List(ctx)
}
