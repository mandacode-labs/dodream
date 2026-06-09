package card

import (
	"context"

	"github.com/microcosm-cc/bluemonday"

	"github.com/mandacode-labs/dodream/internal/errs"
)

// Repository defines the interface for card persistence operations.
type Repository interface {
	Create(ctx context.Context, c *Card) (*Card, error)
	GetByID(ctx context.Context, id ID) (*Card, error)
	Update(ctx context.Context, c *Card) (*Card, error)
	Delete(ctx context.Context, id ID) error
	List(ctx context.Context) ([]*Card, error)
}

// Service provides card-related business logic.
type Service struct {
	repo Repository
}

// NewService creates a new Service with the given repository.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// sanitizeMarkdown removes potentially dangerous HTML from markdown content.
func sanitizeMarkdown(input string) string {
	p := bluemonday.UGCPolicy()
	return p.Sanitize(input)
}

// Create creates a new card.
func (s *Service) Create(ctx context.Context, question string, hint string, content string, creator string) (*Card, error) {
	if question == "" || content == "" {
		return nil, errs.New(errs.ErrInvalidInput, "question and content are required")
	}
	c := New(sanitizeMarkdown(question), sanitizeMarkdown(hint), sanitizeMarkdown(content), creator)
	return s.repo.Create(ctx, c)
}

// GetByID retrieves a card by ID.
func (s *Service) GetByID(ctx context.Context, id ID) (*Card, error) {
	return s.repo.GetByID(ctx, id)
}

// Update modifies an existing card.
func (s *Service) Update(ctx context.Context, id ID, question string, hint string, content string) (*Card, error) {
	if question == "" || content == "" {
		return nil, errs.New(errs.ErrInvalidInput, "question and content are required")
	}
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	c.SetQuestion(sanitizeMarkdown(question))
	c.SetHint(sanitizeMarkdown(hint))
	c.SetContent(sanitizeMarkdown(content))
	return s.repo.Update(ctx, c)
}

// Delete removes a card by ID.
func (s *Service) Delete(ctx context.Context, id ID) error {
	return s.repo.Delete(ctx, id)
}

// List retrieves all cards.
func (s *Service) List(ctx context.Context) ([]*Card, error) {
	return s.repo.List(ctx)
}
