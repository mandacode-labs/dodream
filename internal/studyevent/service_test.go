package studyevent_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mandacode-labs/dodream/internal/errs"
	"github.com/mandacode-labs/dodream/internal/studyevent"
	"github.com/mandacode-labs/dodream/internal/studyevent/mocks"
)

func TestService_CreateEvent(t *testing.T) {
	tests := []struct {
		name        string
		eventType   studyevent.EventType
		quality     *int
		mockSetup   func(repo *mocks.MockRepository)
		wantErr     bool
		errType     errs.ErrorType
		wantPublish bool
	}{
		{
			name:      "valid review event",
			eventType: studyevent.TypeReview,
			quality:   intPtr(4),
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*studyevent.Event")).
					Return(&studyevent.Event{}, nil)
			},
			wantErr:     false,
			wantPublish: false, // no publisher in test
		},
		{
			name:      "invalid quality too high",
			eventType: studyevent.TypeReview,
			quality:   intPtr(6),
			mockSetup: func(repo *mocks.MockRepository) {},
			wantErr:   true,
			errType:   errs.ErrInvalidInput,
		},
		{
			name:      "invalid quality negative",
			eventType: studyevent.TypeReview,
			quality:   intPtr(-1),
			mockSetup: func(repo *mocks.MockRepository) {},
			wantErr:   true,
			errType:   errs.ErrInvalidInput,
		},
		{
			name:      "valid hint shown event without quality",
			eventType: studyevent.TypeHintShown,
			quality:   nil,
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*studyevent.Event")).
					Return(&studyevent.Event{}, nil)
			},
			wantErr: false,
		},
		{
			name:      "repository error",
			eventType: studyevent.TypeReview,
			quality:   intPtr(3),
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*studyevent.Event")).
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

			svc := studyevent.NewService(repo, nil)
			_, err := svc.CreateEvent(
				context.Background(),
				"collection-123",
				"card-123",
				"user-123",
				tt.eventType,
				tt.quality,
				durationPtr(5*time.Second),
				nil,
			)

			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, errs.Is(err, tt.errType))
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestService_GetByID(t *testing.T) {
	tests := []struct {
		name      string
		id        studyevent.ID
		mockSetup func(repo *mocks.MockRepository)
		wantErr   bool
		errType   errs.ErrorType
	}{
		{
			name: "event found",
			id:   studyevent.ID("event-123"),
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetByID(mock.Anything, studyevent.ID("event-123")).
					Return(&studyevent.Event{}, nil)
			},
			wantErr: false,
		},
		{
			name: "event not found",
			id:   studyevent.ID("event-123"),
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetByID(mock.Anything, studyevent.ID("event-123")).
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

			svc := studyevent.NewService(repo, nil)
			_, err := svc.GetByID(context.Background(), tt.id)

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

func durationPtr(d time.Duration) *time.Duration {
	return &d
}
