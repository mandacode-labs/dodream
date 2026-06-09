package collection

import (
	"context"

	"github.com/mandacode-labs/dodream/internal/errs"
)

// Repository defines the storage operations for collections.
type Repository interface {
	Create(ctx context.Context, c *Collection) (*Collection, error)
	GetByID(ctx context.Context, id ID) (*Collection, error)
	Update(ctx context.Context, c *Collection) (*Collection, error)
	Delete(ctx context.Context, id ID) error
	List(ctx context.Context) ([]*Collection, error)
	ListByCreator(ctx context.Context, userID string) ([]*Collection, error)
	AddCard(ctx context.Context, cc *CollectionCard) (*CollectionCard, error)
	RemoveCard(ctx context.Context, collectionCardID CollectionCardID) error
	ListCards(ctx context.Context, collectionID ID) ([]*CollectionCard, error)
}

// Service provides collection-related business logic.
type Service struct {
	repo Repository
}

// NewService creates a new Service with the given repository.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new collection.
func (s *Service) Create(ctx context.Context, name string, creator string) (*Collection, error) {
	if name == "" {
		return nil, errs.New(errs.ErrInvalidInput, "name is required")
	}
	c := New(name, creator)
	return s.repo.Create(ctx, c)
}

// GetByID retrieves a collection by ID.
func (s *Service) GetByID(ctx context.Context, id ID) (*Collection, error) {
	return s.repo.GetByID(ctx, id)
}

// Update modifies an existing collection.
func (s *Service) Update(ctx context.Context, id ID, name string) (*Collection, error) {
	if name == "" {
		return nil, errs.New(errs.ErrInvalidInput, "name is required")
	}
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	c.SetName(name)
	return s.repo.Update(ctx, c)
}

// Delete removes a collection by ID.
func (s *Service) Delete(ctx context.Context, id ID) error {
	return s.repo.Delete(ctx, id)
}

// List retrieves all collections.
func (s *Service) List(ctx context.Context) ([]*Collection, error) {
	return s.repo.List(ctx)
}

// ListByCreator retrieves all collections created by a specific user.
func (s *Service) ListByCreator(ctx context.Context, userID string) ([]*Collection, error) {
	return s.repo.ListByCreator(ctx, userID)
}

// AddCard adds a card to a collection.
func (s *Service) AddCard(ctx context.Context, collectionID ID, cardID string, deckID *string) (*CollectionCard, error) {
	cc := NewCollectionCard(collectionID, cardID, deckID)
	return s.repo.AddCard(ctx, cc)
}

// RemoveCard removes a card from a collection.
func (s *Service) RemoveCard(ctx context.Context, collectionCardID CollectionCardID) error {
	return s.repo.RemoveCard(ctx, collectionCardID)
}

// ListCards retrieves all cards within a specific collection.
func (s *Service) ListCards(ctx context.Context, collectionID ID) ([]*CollectionCard, error) {
	return s.repo.ListCards(ctx, collectionID)
}
