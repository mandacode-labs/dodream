package state

import (
	"context"
	"fmt"
	"time"

	"github.com/mandacode-labs/dodream/ent"
	entcard "github.com/mandacode-labs/dodream/ent/card"
	"github.com/mandacode-labs/dodream/ent/state"
	"github.com/mandacode-labs/dodream/ent/user"
	"github.com/mandacode-labs/dodream/internal/errs"
)

// Store provides database operations for card states.
type Store struct {
	client *ent.Client
}

// NewStore creates a new Store with the given ent client.
func NewStore(client *ent.Client) *Store {
	return &Store{client: client}
}

func mapEntError(op string, err error) error {
	if ent.IsNotFound(err) {
		return errs.Wrap(errs.ErrNotFound, "not found", err)
	}
	if ent.IsConstraintError(err) {
		return errs.Wrap(errs.ErrConflict, "conflict", err)
	}
	return errs.Wrap(errs.ErrInternal, op, err)
}

// Create inserts a new state into the database.
func (s *Store) Create(ctx context.Context, st *State) (*State, error) {
	created, err := s.client.State.Create().
		SetID(st.ID().String()).
		SetNextReviewAt(st.NextReviewAt()).
		SetInterval(st.Interval()).
		SetEaseFactor(st.EaseFactor()).
		SetTotalReviews(st.TotalReviews()).
		SetTotalSuccessful(st.TotalSuccessful()).
		SetStreak(st.Streak()).
		SetReviewsLast1d(st.ReviewsLast1d()).
		SetReviewsLast3d(st.ReviewsLast3d()).
		SetReviewsLast7d(st.ReviewsLast7d()).
		SetUserID(st.UserID()).
		SetCardID(st.CardID()).
		Save(ctx)
	if err != nil {
		return nil, mapEntError("create state", err)
	}

	return s.toCore(created), nil
}

// GetByUserAndCard retrieves a state by user ID and card ID.
func (s *Store) GetByUserAndCard(ctx context.Context, userID string, cardID string) (*State, error) {
	st, err := s.client.State.Query().
		Where(state.And(
			state.HasUserWith(user.ID(userID)),
			state.HasCardWith(entcard.ID(cardID)),
		)).
		WithUser().
		WithCard().
		Only(ctx)
	if err != nil {
		return nil, mapEntError("get state by user and card", err)
	}

	return s.toCore(st), nil
}

// Update modifies an existing state.
func (s *Store) Update(ctx context.Context, st *State) (*State, error) {
	updated, err := s.client.State.UpdateOneID(st.ID().String()).
		SetNextReviewAt(st.NextReviewAt()).
		SetInterval(st.Interval()).
		SetEaseFactor(st.EaseFactor()).
		SetTotalReviews(st.TotalReviews()).
		SetTotalSuccessful(st.TotalSuccessful()).
		SetStreak(st.Streak()).
		SetReviewsLast1d(st.ReviewsLast1d()).
		SetReviewsLast3d(st.ReviewsLast3d()).
		SetReviewsLast7d(st.ReviewsLast7d()).
		SetNillableAvgResponseTimeMs(float64PtrToFloat64(st.AvgResponseTime())).
		SetNillableLastReviewAt(timePtrToTime(st.LastReviewAt())).
		Save(ctx)
	if err != nil {
		return nil, mapEntError("update state", err)
	}

	return s.toCore(updated), nil
}

// Upsert creates or updates a state for a user-card pair.
func (s *Store) Upsert(ctx context.Context, st *State) (*State, error) {
	existing, err := s.GetByUserAndCard(ctx, st.UserID(), st.CardID())
	if err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return s.Create(ctx, st)
		}
		return nil, fmt.Errorf("upsert state: %w", err)
	}

	st.SetID(existing.ID())
	st.SetCreatedAt(existing.CreatedAt())
	return s.Update(ctx, st)
}

// toCore converts an ent State to a domain State.
func (s *Store) toCore(st *ent.State) *State {
	state := New(
		st.Edges.User.ID,
		st.Edges.Card.ID,
	)

	state.SetID(ID(st.ID))
	state.SetCreatedAt(st.CreatedAt)
	state.SetNextReviewAt(st.NextReviewAt)
	state.SetInterval(st.Interval)
	state.SetEaseFactor(st.EaseFactor)

	for i := 0; i < st.TotalReviews; i++ {
		state.IncrementTotalReviews()
	}
	for i := 0; i < st.TotalSuccessful; i++ {
		state.IncrementTotalSuccessful()
	}
	state.SetStreak(st.Streak)
	state.SetReviewsLast1d(st.ReviewsLast1d)
	state.SetReviewsLast3d(st.ReviewsLast3d)
	state.SetReviewsLast7d(st.ReviewsLast7d)

	if st.AvgResponseTimeMs != nil {
		d := time.Duration(*st.AvgResponseTimeMs * float64(time.Millisecond))
		state.SetAvgResponseTime(d)
	}
	if st.LastReviewAt != nil {
		state.SetLastReviewAt(*st.LastReviewAt)
	}

	return state
}

func float64PtrToFloat64(d *time.Duration) *float64 {
	if d == nil {
		return nil
	}
	v := float64(*d) / float64(time.Millisecond)
	return &v
}

func timePtrToTime(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	return t
}
