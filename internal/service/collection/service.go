package collection

import (
	"context"

	"github.com/mandacode-labs/dodream/internal/core"
	"github.com/mandacode-labs/dodream/internal/store"
)

// Service provides collection-related business logic.
type Service struct {
	store *store.CollectionStore
}

// NewService creates a new Service with the given store.
func NewService(store *store.CollectionStore) *Service {
	return &Service{store: store}
}

// Create creates a new collection.
func (s *Service) Create(ctx context.Context, name string, creator core.UserID) (*core.Collection, error) {
	collection := core.NewCollection(name, creator)
	return s.store.Create(ctx, collection)
}

// GetByID retrieves a collection by ID.
func (s *Service) GetByID(ctx context.Context, id core.CollectionID) (*core.Collection, error) {
	return s.store.GetByID(ctx, id)
}

// List retrieves all collections.
func (s *Service) List(ctx context.Context) ([]*core.Collection, error) {
	return s.store.List(ctx)
}

// ListByCreator retrieves all collections created by a specific user.
func (s *Service) ListByCreator(ctx context.Context, userID core.UserID) ([]*core.Collection, error) {
	return s.store.ListByCreator(ctx, userID)
}

// AddCard adds a card to a collection.
func (s *Service) AddCard(ctx context.Context, collectionID core.CollectionID, cardID core.CardID, deckID *core.DeckID) (*core.CollectionCard, error) {
	return s.store.AddCard(ctx, collectionID, cardID, deckID)
}

// RemoveCard removes a card from a collection.
func (s *Service) RemoveCard(ctx context.Context, collectionCardID core.CollectionCardID) error {
	return s.store.RemoveCard(ctx, collectionCardID)
}

// ListCards retrieves all cards within a specific collection.
func (s *Service) ListCards(ctx context.Context, collectionID core.CollectionID) ([]*core.CollectionCard, error) {
	return s.store.ListCards(ctx, collectionID)
}
