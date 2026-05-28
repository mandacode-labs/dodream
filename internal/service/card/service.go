package card

import (
	"context"

	"github.com/mandacode-labs/dodream/internal/core"
	"github.com/mandacode-labs/dodream/internal/store"
)

// Service defines the interface for card-related business logic.
type Service interface {
	Create(ctx context.Context, hint string, content string, creator core.UserID) (*core.Card, error)
	GetByID(ctx context.Context, id core.CardID) (*core.Card, error)
	List(ctx context.Context) ([]*core.Card, error)
}

// CardService implements the Service interface.
type CardService struct {
	store *store.CardStore
}

// NewCardService creates a new CardService with the given store.
func NewCardService(store *store.CardStore) *CardService {
	return &CardService{store: store}
}

// Create creates a new card.
func (s *CardService) Create(ctx context.Context, hint string, content string, creator core.UserID) (*core.Card, error) {
	card := core.NewCard(hint, content, creator)
	return s.store.Create(ctx, card)
}

// GetByID retrieves a card by ID.
func (s *CardService) GetByID(ctx context.Context, id core.CardID) (*core.Card, error) {
	return s.store.GetByID(ctx, id)
}

// List retrieves all cards.
func (s *CardService) List(ctx context.Context) ([]*core.Card, error) {
	return s.store.List(ctx)
}
