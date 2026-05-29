package card

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mandacode-labs/dodream/internal/errs"
)

func TestService_Create(t *testing.T) {
	tests := []struct {
		name      string
		question  string
		hint      string
		content   string
		creator   string
		mockSetup func(repo *MockRepository)
		wantErr   bool
		errType   errs.ErrorType
	}{
		{
			name:     "valid card",
			question: "What is Go?",
			hint:     "A programming language",
			content:  "Go is a statically typed, compiled programming language.",
			creator:  "user-123",
			mockSetup: func(repo *MockRepository) {
				repo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*card.Card")).
					Return(&Card{}, nil)
			},
			wantErr: false,
		},
		{
			name:      "empty question",
			question:  "",
			hint:      "hint",
			content:   "content",
			creator:   "user-123",
			mockSetup: func(repo *MockRepository) {},
			wantErr:   true,
			errType:   errs.ErrInvalidInput,
		},
		{
			name:      "empty content",
			question:  "question",
			hint:      "hint",
			content:   "",
			creator:   "user-123",
			mockSetup: func(repo *MockRepository) {},
			wantErr:   true,
			errType:   errs.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository(t)
			tt.mockSetup(repo)

			svc := NewService(repo)
			_, err := svc.Create(context.Background(), tt.question, tt.hint, tt.content, tt.creator)

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
		id        ID
		mockSetup func(repo *MockRepository)
		wantErr   bool
		errType   errs.ErrorType
	}{
		{
			name: "card found",
			id:   ID("card-123"),
			mockSetup: func(repo *MockRepository) {
				repo.EXPECT().GetByID(mock.Anything, ID("card-123")).
					Return(&Card{}, nil)
			},
			wantErr: false,
		},
		{
			name: "card not found",
			id:   ID("card-123"),
			mockSetup: func(repo *MockRepository) {
				repo.EXPECT().GetByID(mock.Anything, ID("card-123")).
					Return(nil, errs.New(errs.ErrNotFound, "not found"))
			},
			wantErr: true,
			errType: errs.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository(t)
			tt.mockSetup(repo)

			svc := NewService(repo)
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
