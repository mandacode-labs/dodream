package api

import (
	"context"

	"github.com/mandacode-labs/dodream/internal/errs"
	"github.com/mandacode-labs/dodream/internal/user"
	"github.com/mandacode-labs/dodream/pkg/oas"
)

// UserService defines the interface for user business logic.
type UserService interface {
	Create(ctx context.Context, nickname string, providerID string) (*user.User, error)
	GetByID(ctx context.Context, id user.ID) (*user.User, error)
	Update(ctx context.Context, id user.ID, nickname string, providerID string) (*user.User, error)
	Delete(ctx context.Context, id user.ID) error
	List(ctx context.Context) ([]*user.User, error)
}

// CreateUser implements createUser operation.
func (h *Handler) CreateUser(ctx context.Context, req *oas.CreateUserRequest) (*oas.User, error) {
	u, err := h.userService.Create(ctx, req.Nickname, req.ProviderID)
	if err != nil {
		return nil, err
	}
	return mapUser(u), nil
}

// GetUser implements getUser operation.
func (h *Handler) GetUser(ctx context.Context, params oas.GetUserParams) (oas.GetUserRes, error) {
	u, err := h.userService.GetByID(ctx, user.ID(params.ID))
	if err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &oas.GetUserNotFound{}, nil
		}
		return nil, err
	}
	return mapUser(u), nil
}

// ListUsers implements listUsers operation.
func (h *Handler) ListUsers(ctx context.Context) ([]oas.User, error) {
	users, err := h.userService.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]oas.User, len(users))
	for i, u := range users {
		result[i] = *mapUser(u)
	}
	return result, nil
}

// UpdateUser implements updateUser operation.
func (h *Handler) UpdateUser(ctx context.Context, req *oas.UpdateUserRequest, params oas.UpdateUserParams) (oas.UpdateUserRes, error) {
	u, err := h.userService.Update(ctx, user.ID(params.ID), req.Nickname, req.ProviderID)
	if err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &oas.UpdateUserNotFound{}, nil
		}
		return nil, err
	}
	return mapUser(u), nil
}

// DeleteUser implements deleteUser operation.
func (h *Handler) DeleteUser(ctx context.Context, params oas.DeleteUserParams) (oas.DeleteUserRes, error) {
	if err := h.userService.Delete(ctx, user.ID(params.ID)); err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &oas.DeleteUserNotFound{}, nil
		}
		return nil, err
	}
	return &oas.DeleteUserNoContent{}, nil
}
