package state_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mandacode-labs/dodream/internal/engine/state"
	"github.com/mandacode-labs/dodream/internal/engine/state/mocks"
	"github.com/mandacode-labs/dodream/internal/errs"
	"github.com/mandacode-labs/dodream/internal/studyevent"
)

func TestCalculator_CalculateNextState_NewCard(t *testing.T) {
	calc := state.NewCalculator()

	// state.New card with no reviews
	st := state.New("user-1", "card-1")
	st.SetNextReviewAt(time.Now().Add(24 * time.Hour))
	st.SetInterval(0)
	st.SetEaseFactor(2.5)

	// One review with quality 4 (Easy)
	events := []*studyevent.Event{
		studyevent.New(
			studyevent.ID("event-1"),
			"collection-1",
			"card-1",
			"user-1",
			studyevent.TypeReview,
			new(4),
			nil,
			nil,
		),
	}

	now := time.Now()
	newState := calc.CalculateNextState(st, events, now)

	assert.NotNil(t, newState)
	assert.Equal(t, "user-1", newState.UserID())
	assert.Equal(t, "card-1", newState.CardID())
	// After an "Easy" review, the next review should be in the future
	assert.True(t, newState.NextReviewAt().After(now))
}

func TestCalculator_CalculateNextState_ReviewCard(t *testing.T) {
	calc := state.NewCalculator()

	// Card with existing reviews
	st := state.New("user-1", "card-1")
	st.SetNextReviewAt(time.Now())
	st.SetInterval(5.0)
	st.SetEaseFactor(2.5)
	lastReview := time.Now().Add(-24 * time.Hour)
	st.SetLastReviewAt(lastReview)

	// Review with quality 1 (Again) - should decrease interval
	events := []*studyevent.Event{
		studyevent.New(
			studyevent.ID("event-1"),
			"collection-1",
			"card-1",
			"user-1",
			studyevent.TypeReview,
			intPtr(1),
			nil,
			nil,
		),
	}

	now := time.Now()
	newState := calc.CalculateNextState(st, events, now)

	assert.NotNil(t, newState)
	// After "Again", interval should typically be shorter
	assert.True(t, newState.Interval() <= 5.0 || newState.NextReviewAt().Before(now.Add(24*time.Hour)))
}

func TestCalculator_CalculateNextState_NonReviewEvents(t *testing.T) {
	calc := state.NewCalculator()

	st := state.New("user-1", "card-1")
	originalDue := time.Now().Add(24 * time.Hour)
	st.SetNextReviewAt(originalDue)
	st.SetInterval(1.0)
	st.SetEaseFactor(2.5)

	// Hint shown event should not affect state
	events := []*studyevent.Event{
		studyevent.New(
			studyevent.ID("event-1"),
			"collection-1",
			"card-1",
			"user-1",
			studyevent.TypeHintShown,
			nil,
			nil,
			nil,
		),
	}

	newState := calc.CalculateNextState(st, events, time.Now())

	assert.NotNil(t, newState)
	// Non-review events should not change the state significantly
	assert.Equal(t, originalDue.Unix(), newState.NextReviewAt().Unix(), deltaSeconds)
}

func TestService_GetByUserAndCard(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		cardID    string
		mockSetup func(repo *mocks.MockRepository)
		wantErr   bool
		errType   errs.ErrorType
	}{
		{
			name:   "state found",
			userID: "user-1",
			cardID: "card-1",
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetByUserAndCard(mock.Anything, "user-1", "card-1").
					Return(&state.State{}, nil)
			},
			wantErr: false,
		},
		{
			name:   "state not found",
			userID: "user-1",
			cardID: "card-1",
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetByUserAndCard(mock.Anything, "user-1", "card-1").
					Return(nil, errs.New(errs.ErrNotFound, "not found"))
			},
			wantErr: true,
			errType: errs.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockRepository(t)
			tt.mockSetup(repo)

			svc := state.NewService(repo)
			_, err := svc.GetByUserAndCard(context.Background(), tt.userID, tt.cardID)

			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, errs.Is(err, tt.errType))
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestService_Upsert(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(repo *mocks.MockRepository)
		wantErr   bool
		errType   errs.ErrorType
	}{
		{
			name: "upsert success",
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().Upsert(mock.Anything, mock.AnythingOfType("*state.State")).
					Return(&state.State{}, nil)
			},
			wantErr: false,
		},
		{
			name: "upsert error",
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().Upsert(mock.Anything, mock.AnythingOfType("*state.State")).
					Return(nil, errs.New(errs.ErrInternal, "db error"))
			},
			wantErr: true,
			errType: errs.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockRepository(t)
			tt.mockSetup(repo)

			svc := state.NewService(repo)
			st := state.New("user-1", "card-1")
			_, err := svc.Upsert(context.Background(), st)

			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, errs.Is(err, tt.errType))
				return
			}
			require.NoError(t, err)
		})
	}
}

func intPtr(i int) *int {
	return &i
}

// deltaSeconds allows small time differences due to execution time
const deltaSeconds = 2
