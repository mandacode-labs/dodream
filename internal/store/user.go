package store

import (
	"context"
	"fmt"

	"github.com/mandacode-labs/dodream/ent"
	"github.com/mandacode-labs/dodream/ent/user"
	"github.com/mandacode-labs/dodream/internal/core"
)

// UserStore provides database operations for users.
type UserStore struct {
	client *ent.Client
}

// NewUserStore creates a new UserStore with the given ent client.
func NewUserStore(client *ent.Client) *UserStore {
	return &UserStore{client: client}
}

// Create inserts a new user into the database.
func (s *UserStore) Create(ctx context.Context, u *core.User) (*core.User, error) {
	created, err := s.client.User.Create().
		SetID(u.ID().String()).
		SetNickname(u.Nickname()).
		SetProviderID(u.ProviderID()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return core.NewUser(
		core.UserID(created.ID),
		created.Nickname,
		created.ProviderID,
	), nil
}

// GetByID retrieves a user by their ID.
func (s *UserStore) GetByID(ctx context.Context, id core.UserID) (*core.User, error) {
	u, err := s.client.User.Query().
		Where(user.ID(id.String())).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return core.NewUser(
		core.UserID(u.ID),
		u.Nickname,
		u.ProviderID,
	), nil
}

// Update modifies an existing user.
func (s *UserStore) Update(ctx context.Context, u *core.User) (*core.User, error) {
	updated, err := s.client.User.UpdateOneID(u.ID().String()).
		SetNickname(u.Nickname()).
		SetProviderID(u.ProviderID()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return core.NewUser(
		core.UserID(updated.ID),
		updated.Nickname,
		updated.ProviderID,
	), nil
}

// Delete removes a user by their ID.
func (s *UserStore) Delete(ctx context.Context, id core.UserID) error {
	err := s.client.User.DeleteOneID(id.String()).Exec(ctx)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

// List retrieves all users.
func (s *UserStore) List(ctx context.Context) ([]*core.User, error) {
	users, err := s.client.User.Query().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	result := make([]*core.User, len(users))
	for i, u := range users {
		result[i] = core.NewUser(
			core.UserID(u.ID),
			u.Nickname,
			u.ProviderID,
		)
	}
	return result, nil
}
