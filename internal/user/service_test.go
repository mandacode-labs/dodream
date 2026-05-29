package user

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
		name       string
		nickname   string
		providerID string
		mockSetup  func(repo *MockRepository)
		wantErr    bool
		errType    errs.ErrorType
	}{
		{
			name:       "valid user",
			nickname:   "testuser",
			providerID: "provider_123",
			mockSetup: func(repo *MockRepository) {
				repo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*user.User")).
					Return(&User{}, nil)
			},
			wantErr: false,
		},
		{
			name:       "empty nickname",
			nickname:   "",
			providerID: "provider_123",
			mockSetup:  func(repo *MockRepository) {},
			wantErr:    true,
			errType:    errs.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository(t)
			tt.mockSetup(repo)

			svc := NewService(repo)
			_, err := svc.Create(context.Background(), tt.nickname, tt.providerID)

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
			name: "user found",
			id:   ID("user-123"),
			mockSetup: func(repo *MockRepository) {
				repo.EXPECT().GetByID(mock.Anything, ID("user-123")).
					Return(&User{}, nil)
			},
			wantErr: false,
		},
		{
			name: "user not found",
			id:   ID("user-123"),
			mockSetup: func(repo *MockRepository) {
				repo.EXPECT().GetByID(mock.Anything, ID("user-123")).
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
