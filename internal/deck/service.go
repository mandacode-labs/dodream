package deck

import (
	"context"

	"github.com/mandacode-labs/dodream/internal/errs"
)

// Repository defines the interface for deck data access.
type Repository interface {
	Create(ctx context.Context, d *Deck) (*Deck, error)
	GetByID(ctx context.Context, id ID) (*Deck, error)
	Update(ctx context.Context, d *Deck) (*Deck, error)
	Delete(ctx context.Context, id ID) error
	List(ctx context.Context) ([]*Deck, error)
	ListByCreator(ctx context.Context, userID string) ([]*Deck, error)
}

// Service provides deck-related business logic.
type Service struct {
	repo Repository
}

// NewService creates a new Service with the given repository.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new deck.
func (s *Service) Create(ctx context.Context, name string, creator string) (*Deck, error) {
	if name == "" {
		return nil, errs.New(errs.ErrInvalidInput, "name cannot be empty")
	}
	d := New(name, creator)
	return s.repo.Create(ctx, d)
}

// GetByID retrieves a deck by ID.
func (s *Service) GetByID(ctx context.Context, id ID) (*Deck, error) {
	return s.repo.GetByID(ctx, id)
}

// Update modifies an existing deck.
func (s *Service) Update(ctx context.Context, id ID, name string) (*Deck, error) {
	if name == "" {
		return nil, errs.New(errs.ErrInvalidInput, "name cannot be empty")
	}
	d, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	d.SetName(name)
	return s.repo.Update(ctx, d)
}

// Delete removes a deck by ID.
func (s *Service) Delete(ctx context.Context, id ID) error {
	return s.repo.Delete(ctx, id)
}

// List retrieves all decks.
func (s *Service) List(ctx context.Context) ([]*Deck, error) {
	return s.repo.List(ctx)
}

// ListByCreator retrieves all decks created by a specific user.
func (s *Service) ListByCreator(ctx context.Context, userID string) ([]*Deck, error) {
	return s.repo.ListByCreator(ctx, userID)
}
