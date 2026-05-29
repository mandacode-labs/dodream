package deck_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mandacode-labs/dodream/internal/deck"
	"github.com/mandacode-labs/dodream/internal/deck/mocks"
	"github.com/mandacode-labs/dodream/internal/errs"
)

func TestService_Create(t *testing.T) {
	tests := []struct {
		name      string
		name_     string
		creator   string
		mockSetup func(repo *mocks.MockRepository)
		wantErr   bool
		errType   errs.ErrorType
	}{
		{
			name:    "valid deck",
			name_:   "My deck.Deck",
			creator: "user-123",
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*deck.Deck")).
					Return(&deck.Deck{}, nil)
			},
			wantErr: false,
		},
		{
			name:      "empty name",
			name_:     "",
			creator:   "user-123",
			mockSetup: func(repo *mocks.MockRepository) {},
			wantErr:   true,
			errType:   errs.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockRepository(t)
			tt.mockSetup(repo)

			svc := deck.NewService(repo)
			_, err := svc.Create(context.Background(), tt.name_, tt.creator)

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
		id        deck.ID
		mockSetup func(repo *mocks.MockRepository)
		wantErr   bool
		errType   errs.ErrorType
	}{
		{
			name: "deck found",
			id:   deck.ID("deck-123"),
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetByID(mock.Anything, deck.ID("deck-123")).
					Return(&deck.Deck{}, nil)
			},
			wantErr: false,
		},
		{
			name: "deck not found",
			id:   deck.ID("deck-123"),
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetByID(mock.Anything, deck.ID("deck-123")).
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

			svc := deck.NewService(repo)
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
