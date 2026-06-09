package user

import (
	"context"

	"github.com/mandacode-labs/dodream/internal/errs"
)

// Repository is the persistence interface consumed by Service.
type Repository interface {
	Create(ctx context.Context, u *User) (*User, error)
	GetByID(ctx context.Context, id ID) (*User, error)
	Update(ctx context.Context, u *User) (*User, error)
	Delete(ctx context.Context, id ID) error
	List(ctx context.Context) ([]*User, error)
}

// Service provides user-related business logic.
type Service struct {
	repo Repository
}

// NewService creates a new Service with the given repository.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new user.
func (s *Service) Create(ctx context.Context, nickname string, providerID string) (*User, error) {
	if nickname == "" {
		return nil, errs.New(errs.ErrInvalidInput, "nickname is required")
	}
	u := New(nickname, providerID)
	return s.repo.Create(ctx, u)
}

// GetByID retrieves a user by ID.
func (s *Service) GetByID(ctx context.Context, id ID) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

// Update modifies an existing user.
func (s *Service) Update(ctx context.Context, id ID, nickname string, providerID string) (*User, error) {
	if nickname == "" {
		return nil, errs.New(errs.ErrInvalidInput, "nickname is required")
	}
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	u.SetNickname(nickname)
	u.SetProviderID(providerID)
	return s.repo.Update(ctx, u)
}

// Delete removes a user by ID.
func (s *Service) Delete(ctx context.Context, id ID) error {
	return s.repo.Delete(ctx, id)
}

// List retrieves all users.
func (s *Service) List(ctx context.Context) ([]*User, error) {
	return s.repo.List(ctx)
}
