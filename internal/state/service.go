package state

import (
	"context"
	"time"

	"github.com/open-spaced-repetition/go-fsrs"

	"github.com/mandacode-labs/dodream/internal/studyevent"
)

// Calculator wraps the FSRS scheduler for SRS calculations.
type Calculator struct {
	params fsrs.Parameters
}

// NewCalculator creates a new FSRS calculator with default parameters.
func NewCalculator() *Calculator {
	return &Calculator{
		params: fsrs.DefaultParam(),
	}
}

// CalculateNextState computes the next SRS state from the current state and study events.
func (c *Calculator) CalculateNextState(st *State, events []*studyevent.Event, now time.Time) *State {
	card := fsrs.Card{
		Due:           st.NextReviewAt(),
		Stability:     st.Interval(),
		Difficulty:    st.EaseFactor(),
		ElapsedDays:   0,
		ScheduledDays: uint64(st.Interval()),
		Reps:          uint64(st.TotalReviews()),
		Lapses:        uint64(st.TotalReviews() - st.TotalSuccessful()),
		State:         fsrs.New,
	}

	if st.LastReviewAt() != nil {
		card.LastReview = *st.LastReviewAt()
	}

	if st.TotalReviews() > 0 {
		card.State = fsrs.Review
	}

	for _, event := range events {
		if event.EventType() != studyevent.TypeReview || event.Quality() == nil {
			continue
		}

		rating := mapQualityToRating(*event.Quality())
		result := c.params.Repeat(card, event.CreatedAt())
		if info, ok := result[rating]; ok {
			card = info.Card
		}
	}

	newState := New(st.UserID(), st.CardID())
	newState.SetNextReviewAt(card.Due)
	newState.SetInterval(float64(card.ScheduledDays))
	newState.SetEaseFactor(card.Difficulty)

	return newState
}

func mapQualityToRating(quality int) fsrs.Rating {
	switch quality {
	case 0, 1:
		return fsrs.Again
	case 2:
		return fsrs.Hard
	case 3:
		return fsrs.Good
	case 4, 5:
		return fsrs.Easy
	default:
		return fsrs.Good
	}
}

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
