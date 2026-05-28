package collection

import (
	"context"

	"github.com/mandacode-labs/dodream/internal/core"
	"github.com/mandacode-labs/dodream/internal/store"
)

// Service defines the interface for collection-related business logic.
type Service interface {
	Create(ctx context.Context, name string, creator core.UserID) (*core.Collection, error)
	GetByID(ctx context.Context, id core.CollectionID) (*core.Collection, error)
	List(ctx context.Context) ([]*core.Collection, error)
	ListByCreator(ctx context.Context, userID core.UserID) ([]*core.Collection, error)
	AddCard(ctx context.Context, collectionID core.CollectionID, cardID core.CardID, deckID *core.DeckID) (*core.CollectionCard, error)
	RemoveCard(ctx context.Context, collectionCardID core.CollectionCardID) error
	ListCards(ctx context.Context, collectionID core.CollectionID) ([]*core.CollectionCard, error)
}

// CollectionService implements the Service interface.
type CollectionService struct {
	store *store.CollectionStore
}

// NewCollectionService creates a new CollectionService with the given store.
func NewCollectionService(store *store.CollectionStore) *CollectionService {
	return &CollectionService{store: store}
}

// Create creates a new collection.
func (s *CollectionService) Create(ctx context.Context, name string, creator core.UserID) (*core.Collection, error) {
	collection := core.NewCollection(name, creator)
	return s.store.Create(ctx, collection)
}

// GetByID retrieves a collection by ID.
func (s *CollectionService) GetByID(ctx context.Context, id core.CollectionID) (*core.Collection, error) {
	return s.store.GetByID(ctx, id)
}

// List retrieves all collections.
func (s *CollectionService) List(ctx context.Context) ([]*core.Collection, error) {
	return s.store.List(ctx)
}

// ListByCreator retrieves all collections created by a specific user.
func (s *CollectionService) ListByCreator(ctx context.Context, userID core.UserID) ([]*core.Collection, error) {
	return s.store.ListByCreator(ctx, userID)
}

// AddCard adds a card to a collection.
func (s *CollectionService) AddCard(ctx context.Context, collectionID core.CollectionID, cardID core.CardID, deckID *core.DeckID) (*core.CollectionCard, error) {
	return s.store.AddCard(ctx, collectionID, cardID, deckID)
}

// RemoveCard removes a card from a collection.
func (s *CollectionService) RemoveCard(ctx context.Context, collectionCardID core.CollectionCardID) error {
	return s.store.RemoveCard(ctx, collectionCardID)
}

// ListCards retrieves all cards within a specific collection.
func (s *CollectionService) ListCards(ctx context.Context, collectionID core.CollectionID) ([]*core.CollectionCard, error) {
	return s.store.ListCards(ctx, collectionID)
}
