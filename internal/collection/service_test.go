package collection_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mandacode-labs/dodream/internal/collection"
	"github.com/mandacode-labs/dodream/internal/collection/mocks"
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
			name:    "valid collection",
			name_:   "My collection.Collection",
			creator: "user-123",
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*collection.Collection")).
					Return(&collection.Collection{}, nil)
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

			svc := collection.NewService(repo)
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
		id        collection.ID
		mockSetup func(repo *mocks.MockRepository)
		wantErr   bool
		errType   errs.ErrorType
	}{
		{
			name: "collection found",
			id:   collection.ID("collection-123"),
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetByID(mock.Anything, collection.ID("collection-123")).
					Return(&collection.Collection{}, nil)
			},
			wantErr: false,
		},
		{
			name: "collection not found",
			id:   collection.ID("collection-123"),
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetByID(mock.Anything, collection.ID("collection-123")).
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

			svc := collection.NewService(repo)
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
