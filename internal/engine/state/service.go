package state

import (
	"context"
)

// Repository defines the data access interface for states.
type Repository interface {
	Create(ctx context.Context, st *State) (*State, error)
	GetByUserAndCard(ctx context.Context, userID string, cardID string) (*State, error)
	Update(ctx context.Context, st *State) (*State, error)
	Upsert(ctx context.Context, st *State) (*State, error)
}

// Service provides state-related business logic.
type Service struct {
	repo Repository
}

// NewService creates a new Service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetByUserAndCard retrieves a state.
func (s *Service) GetByUserAndCard(ctx context.Context, userID string, cardID string) (*State, error) {
	return s.repo.GetByUserAndCard(ctx, userID, cardID)
}

// Upsert creates or updates a state.
func (s *Service) Upsert(ctx context.Context, st *State) (*State, error) {
	return s.repo.Upsert(ctx, st)
}
