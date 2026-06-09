package user

import (
	"context"

	"github.com/mandacode-labs/dodream/ent"
	entuser "github.com/mandacode-labs/dodream/ent/user"
	"github.com/mandacode-labs/dodream/internal/errs"
)

// Store provides database operations for users.
type Store struct {
	client *ent.Client
}

// NewStore creates a new Store with the given ent client.
func NewStore(client *ent.Client) *Store {
	return &Store{client: client}
}

// Compile-time check that Store implements Repository.
var _ Repository = (*Store)(nil)

func mapEntError(op string, err error) error {
	if ent.IsNotFound(err) {
		return errs.Wrap(errs.ErrNotFound, op, err)
	}
	if ent.IsConstraintError(err) {
		return errs.Wrap(errs.ErrConflict, op, err)
	}
	return errs.Wrap(errs.ErrInternal, op, err)
}

// Create inserts a new user into the database.
func (s *Store) Create(ctx context.Context, u *User) (*User, error) {
	created, err := s.client.User.Create().
		SetID(u.ID().String()).
		SetNickname(u.Nickname()).
		SetProviderID(u.ProviderID()).
		Save(ctx)
	if err != nil {
		return nil, mapEntError("create user", err)
	}
	return NewWithID(
		ID(created.ID),
		created.Nickname,
		created.ProviderID,
		created.CreatedAt,
		created.UpdatedAt,
	), nil
}

// GetByID retrieves a user by their ID.
func (s *Store) GetByID(ctx context.Context, id ID) (*User, error) {
	u, err := s.client.User.Query().
		Where(entuser.ID(id.String())).
		Only(ctx)
	if err != nil {
		return nil, mapEntError("get user by id", err)
	}
	return NewWithID(
		ID(u.ID),
		u.Nickname,
		u.ProviderID,
		u.CreatedAt,
		u.UpdatedAt,
	), nil
}

// Update modifies an existing user.
func (s *Store) Update(ctx context.Context, u *User) (*User, error) {
	updated, err := s.client.User.UpdateOneID(u.ID().String()).
		SetNickname(u.Nickname()).
		SetProviderID(u.ProviderID()).
		Save(ctx)
	if err != nil {
		return nil, mapEntError("update user", err)
	}
	return NewWithID(
		ID(updated.ID),
		updated.Nickname,
		updated.ProviderID,
		updated.CreatedAt,
		updated.UpdatedAt,
	), nil
}

// Delete removes a user by their ID.
func (s *Store) Delete(ctx context.Context, id ID) error {
	err := s.client.User.DeleteOneID(id.String()).Exec(ctx)
	if err != nil {
		return mapEntError("delete user", err)
	}
	return nil
}

// List retrieves all users.
func (s *Store) List(ctx context.Context) ([]*User, error) {
	users, err := s.client.User.Query().All(ctx)
	if err != nil {
		return nil, mapEntError("list users", err)
	}

	result := make([]*User, len(users))
	for i, u := range users {
		result[i] = NewWithID(
			ID(u.ID),
			u.Nickname,
			u.ProviderID,
			u.CreatedAt,
			u.UpdatedAt,
		)
	}
	return result, nil
}
