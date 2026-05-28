package user

import (
	"context"

	"github.com/mandacode-labs/dodream/internal/core"
	"github.com/mandacode-labs/dodream/internal/store"
)

// Service provides user-related business logic.
type Service struct {
	store *store.UserStore
}

// NewService creates a new Service with the given store.
func NewService(store *store.UserStore) *Service {
	return &Service{store: store}
}

// Create creates a new user.
func (s *Service) Create(ctx context.Context, nickname string, providerID string) (*core.User, error) {
	user := core.NewUser(nickname, providerID)
	return s.store.Create(ctx, user)
}

// GetByID retrieves a user by ID.
func (s *Service) GetByID(ctx context.Context, id core.UserID) (*core.User, error) {
	return s.store.GetByID(ctx, id)
}

// List retrieves all users.
func (s *Service) List(ctx context.Context) ([]*core.User, error) {
	return s.store.List(ctx)
}
