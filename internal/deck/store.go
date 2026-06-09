package deck

import (
	"context"

	"github.com/mandacode-labs/dodream/ent"
	entdeck "github.com/mandacode-labs/dodream/ent/deck"
	"github.com/mandacode-labs/dodream/ent/user"
	"github.com/mandacode-labs/dodream/internal/errs"
)

// Store provides database operations for decks.
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

// Create inserts a new deck into the database.
func (s *Store) Create(ctx context.Context, d *Deck) (*Deck, error) {
	created, err := s.client.Deck.Create().
		SetID(d.ID().String()).
		SetName(d.Name()).
		SetCreatorID(d.CreatorID()).
		Save(ctx)
	if err != nil {
		return nil, mapEntError("create deck", err)
	}
	return NewWithID(
		ID(created.ID),
		created.Name,
		d.CreatorID(),
		created.CreatedAt,
		created.UpdatedAt,
	), nil
}

// GetByID retrieves a deck by ID.
func (s *Store) GetByID(ctx context.Context, id ID) (*Deck, error) {
	d, err := s.client.Deck.Query().
		Where(entdeck.ID(id.String())).
		WithCreator().
		Only(ctx)
	if err != nil {
		return nil, mapEntError("get deck by id", err)
	}
	return NewWithID(
		ID(d.ID),
		d.Name,
		d.Edges.Creator.ID,
		d.CreatedAt,
		d.UpdatedAt,
	), nil
}

// Update modifies an existing deck.
func (s *Store) Update(ctx context.Context, d *Deck) (*Deck, error) {
	updated, err := s.client.Deck.UpdateOneID(d.ID().String()).
		SetName(d.Name()).
		Save(ctx)
	if err != nil {
		return nil, mapEntError("update deck", err)
	}
	return NewWithID(
		ID(updated.ID),
		updated.Name,
		d.CreatorID(),
		updated.CreatedAt,
		updated.UpdatedAt,
	), nil
}

// Delete removes a deck by ID.
func (s *Store) Delete(ctx context.Context, id ID) error {
	err := s.client.Deck.DeleteOneID(id.String()).Exec(ctx)
	if err != nil {
		return mapEntError("delete deck", err)
	}
	return nil
}

// List retrieves all decks.
func (s *Store) List(ctx context.Context) ([]*Deck, error) {
	decks, err := s.client.Deck.Query().WithCreator().All(ctx)
	if err != nil {
		return nil, mapEntError("list decks", err)
	}

	result := make([]*Deck, len(decks))
	for i, d := range decks {
		result[i] = NewWithID(
			ID(d.ID),
			d.Name,
			d.Edges.Creator.ID,
			d.CreatedAt,
			d.UpdatedAt,
		)
	}
	return result, nil
}

// ListByCreator retrieves all decks created by a specific user.
func (s *Store) ListByCreator(ctx context.Context, userID string) ([]*Deck, error) {
	decks, err := s.client.Deck.Query().
		Where(entdeck.HasCreatorWith(user.ID(userID))).
		WithCreator().
		All(ctx)
	if err != nil {
		return nil, mapEntError("list decks by creator", err)
	}

	result := make([]*Deck, len(decks))
	for i, d := range decks {
		result[i] = NewWithID(
			ID(d.ID),
			d.Name,
			d.Edges.Creator.ID,
			d.CreatedAt,
			d.UpdatedAt,
		)
	}
	return result, nil
}
