package deck

import (
	"context"

	"github.com/mandacode-labs/dodream/internal/core"
	"github.com/mandacode-labs/dodream/internal/store"
)

// Service defines the interface for deck-related business logic.
type Service interface {
	Create(ctx context.Context, name string, creator core.UserID) (*core.Deck, error)
	GetByID(ctx context.Context, id core.DeckID) (*core.Deck, error)
	List(ctx context.Context) ([]*core.Deck, error)
	ListByCreator(ctx context.Context, userID core.UserID) ([]*core.Deck, error)
}

// DeckService implements the Service interface.
type DeckService struct {
	store *store.DeckStore
}

// NewDeckService creates a new DeckService with the given store.
func NewDeckService(store *store.DeckStore) *DeckService {
	return &DeckService{store: store}
}

// Create creates a new deck.
func (s *DeckService) Create(ctx context.Context, name string, creator core.UserID) (*core.Deck, error) {
	deck := core.NewDeck(name, creator)
	return s.store.Create(ctx, deck)
}

// GetByID retrieves a deck by ID.
func (s *DeckService) GetByID(ctx context.Context, id core.DeckID) (*core.Deck, error) {
	return s.store.GetByID(ctx, id)
}

// List retrieves all decks.
func (s *DeckService) List(ctx context.Context) ([]*core.Deck, error) {
	return s.store.List(ctx)
}

// ListByCreator retrieves all decks created by a specific user.
func (s *DeckService) ListByCreator(ctx context.Context, userID core.UserID) ([]*core.Deck, error) {
	return s.store.ListByCreator(ctx, userID)
}
