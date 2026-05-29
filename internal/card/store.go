package card

import (
	"context"

	"github.com/mandacode-labs/dodream/ent"
	entcard "github.com/mandacode-labs/dodream/ent/card"
	"github.com/mandacode-labs/dodream/internal/errs"
)

// Store provides database operations for cards.
type Store struct {
	client *ent.Client
}

// NewStore creates a new Store with the given ent client.
func NewStore(client *ent.Client) *Store {
	return &Store{client: client}
}

func mapEntError(op string, err error) error {
	if ent.IsNotFound(err) {
		return errs.Wrap(errs.ErrNotFound, op, err)
	}
	if ent.IsConstraintError(err) {
		return errs.Wrap(errs.ErrConflict, op, err)
	}
	return errs.Wrap(errs.ErrInternal, op, err)
}

// Create inserts a new card into the database.
func (s *Store) Create(ctx context.Context, c *Card) (*Card, error) {
	created, err := s.client.Card.Create().
		SetID(c.ID().String()).
		SetQuestion(c.Question()).
		SetHint(c.Hint()).
		SetContent(c.Content()).
		SetCreatorID(c.Creator()).
		Save(ctx)
	if err != nil {
		return nil, mapEntError("create card", err)
	}
	return NewWithID(
		ID(created.ID),
		created.Question,
		created.Hint,
		created.Content,
		c.Creator(),
		created.CreatedAt,
		created.UpdatedAt,
	), nil
}

// GetByID retrieves a card by ID.
func (s *Store) GetByID(ctx context.Context, id ID) (*Card, error) {
	c, err := s.client.Card.Query().
		Where(entcard.ID(id.String())).
		WithCreator().
		Only(ctx)
	if err != nil {
		return nil, mapEntError("get card by id", err)
	}
	return NewWithID(
		ID(c.ID),
		c.Question,
		c.Hint,
		c.Content,
		c.Edges.Creator.ID,
		c.CreatedAt,
		c.UpdatedAt,
	), nil
}

// Update modifies an existing card.
func (s *Store) Update(ctx context.Context, c *Card) (*Card, error) {
	updated, err := s.client.Card.UpdateOneID(c.ID().String()).
		SetQuestion(c.Question()).
		SetHint(c.Hint()).
		SetContent(c.Content()).
		Save(ctx)
	if err != nil {
		return nil, mapEntError("update card", err)
	}
	return NewWithID(
		ID(updated.ID),
		updated.Question,
		updated.Hint,
		updated.Content,
		c.Creator(),
		updated.CreatedAt,
		updated.UpdatedAt,
	), nil
}

// Delete removes a card by ID.
func (s *Store) Delete(ctx context.Context, id ID) error {
	err := s.client.Card.DeleteOneID(id.String()).Exec(ctx)
	if err != nil {
		return mapEntError("delete card", err)
	}
	return nil
}

// List retrieves all cards.
func (s *Store) List(ctx context.Context) ([]*Card, error) {
	cards, err := s.client.Card.Query().WithCreator().All(ctx)
	if err != nil {
		return nil, mapEntError("list cards", err)
	}

	result := make([]*Card, len(cards))
	for i, c := range cards {
		result[i] = NewWithID(
			ID(c.ID),
			c.Question,
			c.Hint,
			c.Content,
			c.Edges.Creator.ID,
			c.CreatedAt,
			c.UpdatedAt,
		)
	}
	return result, nil
}
