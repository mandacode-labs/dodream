package user

import (
	"context"

	"github.com/mandacode-labs/dodream/internal/core"
	"github.com/mandacode-labs/dodream/internal/store"
)

// Service defines the interface for user-related business logic.
type Service interface {
	Create(ctx context.Context, nickname string, providerID string) (*core.User, error)
	GetByID(ctx context.Context, id core.UserID) (*core.User, error)
	List(ctx context.Context) ([]*core.User, error)
}

// UserService implements the Service interface.
type UserService struct {
	store *store.UserStore
}

// NewUserService creates a new UserService with the given store.
func NewUserService(store *store.UserStore) *UserService {
	return &UserService{store: store}
}

// Create creates a new user.
func (s *UserService) Create(ctx context.Context, nickname string, providerID string) (*core.User, error) {
	user := core.NewUser(nickname, providerID)
	return s.store.Create(ctx, user)
}

// GetByID retrieves a user by ID.
func (s *UserService) GetByID(ctx context.Context, id core.UserID) (*core.User, error) {
	return s.store.GetByID(ctx, id)
}

// List retrieves all users.
func (s *UserService) List(ctx context.Context) ([]*core.User, error) {
	return s.store.List(ctx)
}
