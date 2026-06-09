package api

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mandacode-labs/dodream/internal/card"
	"github.com/mandacode-labs/dodream/internal/collection"
	"github.com/mandacode-labs/dodream/internal/deck"
	"github.com/mandacode-labs/dodream/internal/errs"
	"github.com/mandacode-labs/dodream/internal/studyevent"
	"github.com/mandacode-labs/dodream/internal/user"
	"github.com/mandacode-labs/dodream/pkg/oas"
)

// Manual mocks for handler service interfaces

type mockUserService struct {
	createFunc func(ctx context.Context, nickname, providerID string) (*user.User, error)
	getFunc    func(ctx context.Context, id user.ID) (*user.User, error)
	updateFunc func(ctx context.Context, id user.ID, nickname, providerID string) (*user.User, error)
	deleteFunc func(ctx context.Context, id user.ID) error
	listFunc   func(ctx context.Context) ([]*user.User, error)
}

func (m *mockUserService) Create(ctx context.Context, nickname, providerID string) (*user.User, error) {
	return m.createFunc(ctx, nickname, providerID)
}

func (m *mockUserService) GetByID(ctx context.Context, id user.ID) (*user.User, error) {
	return m.getFunc(ctx, id)
}

func (m *mockUserService) Update(ctx context.Context, id user.ID, nickname, providerID string) (*user.User, error) {
	return m.updateFunc(ctx, id, nickname, providerID)
}

func (m *mockUserService) Delete(ctx context.Context, id user.ID) error {
	return m.deleteFunc(ctx, id)
}

func (m *mockUserService) List(ctx context.Context) ([]*user.User, error) {
	return m.listFunc(ctx)
}

type mockCardService struct {
	createFunc func(ctx context.Context, question, hint, content, creator string) (*card.Card, error)
	getFunc    func(ctx context.Context, id card.ID) (*card.Card, error)
	updateFunc func(ctx context.Context, id card.ID, question, hint, content string) (*card.Card, error)
	deleteFunc func(ctx context.Context, id card.ID) error
	listFunc   func(ctx context.Context) ([]*card.Card, error)
}

func (m *mockCardService) Create(ctx context.Context, question, hint, content, creator string) (*card.Card, error) {
	return m.createFunc(ctx, question, hint, content, creator)
}

func (m *mockCardService) GetByID(ctx context.Context, id card.ID) (*card.Card, error) {
	return m.getFunc(ctx, id)
}

func (m *mockCardService) Update(ctx context.Context, id card.ID, question, hint, content string) (*card.Card, error) {
	return m.updateFunc(ctx, id, question, hint, content)
}

func (m *mockCardService) Delete(ctx context.Context, id card.ID) error {
	return m.deleteFunc(ctx, id)
}

func (m *mockCardService) List(ctx context.Context) ([]*card.Card, error) {
	return m.listFunc(ctx)
}

type mockDeckService struct {
	createFunc func(ctx context.Context, name, creator string) (*deck.Deck, error)
	getFunc    func(ctx context.Context, id deck.ID) (*deck.Deck, error)
	updateFunc func(ctx context.Context, id deck.ID, name string) (*deck.Deck, error)
	deleteFunc func(ctx context.Context, id deck.ID) error
	listFunc   func(ctx context.Context) ([]*deck.Deck, error)
}

func (m *mockDeckService) Create(ctx context.Context, name, creator string) (*deck.Deck, error) {
	return m.createFunc(ctx, name, creator)
}

func (m *mockDeckService) GetByID(ctx context.Context, id deck.ID) (*deck.Deck, error) {
	return m.getFunc(ctx, id)
}

func (m *mockDeckService) Update(ctx context.Context, id deck.ID, name string) (*deck.Deck, error) {
	return m.updateFunc(ctx, id, name)
}

func (m *mockDeckService) Delete(ctx context.Context, id deck.ID) error {
	return m.deleteFunc(ctx, id)
}

func (m *mockDeckService) List(ctx context.Context) ([]*deck.Deck, error) {
	return m.listFunc(ctx)
}

type mockCollectionService struct {
	createFunc func(ctx context.Context, name, creator string) (*collection.Collection, error)
	getFunc    func(ctx context.Context, id collection.ID) (*collection.Collection, error)
	updateFunc func(ctx context.Context, id collection.ID, name string) (*collection.Collection, error)
	deleteFunc func(ctx context.Context, id collection.ID) error
	listFunc   func(ctx context.Context) ([]*collection.Collection, error)
}

func (m *mockCollectionService) Create(ctx context.Context, name, creator string) (*collection.Collection, error) {
	return m.createFunc(ctx, name, creator)
}

func (m *mockCollectionService) GetByID(ctx context.Context, id collection.ID) (*collection.Collection, error) {
	return m.getFunc(ctx, id)
}

func (m *mockCollectionService) Update(ctx context.Context, id collection.ID, name string) (*collection.Collection, error) {
	return m.updateFunc(ctx, id, name)
}

func (m *mockCollectionService) Delete(ctx context.Context, id collection.ID) error {
	return m.deleteFunc(ctx, id)
}

func (m *mockCollectionService) List(ctx context.Context) ([]*collection.Collection, error) {
	return m.listFunc(ctx)
}

type mockStudyEventService struct {
	createFunc     func(ctx context.Context, collectionID, cardID, userID string, eventType studyevent.EventType, quality *int, responseTime *time.Duration, previousDeckID *string) (*studyevent.Event, error)
	getByIDFunc    func(ctx context.Context, id studyevent.ID) (*studyevent.Event, error)
	listRecentFunc func(ctx context.Context, userID, cardID string, limit int) ([]*studyevent.Event, error)
}

func (m *mockStudyEventService) CreateEvent(ctx context.Context, collectionID, cardID, userID string, eventType studyevent.EventType, quality *int, responseTime *time.Duration, previousDeckID *string) (*studyevent.Event, error) {
	return m.createFunc(ctx, collectionID, cardID, userID, eventType, quality, responseTime, previousDeckID)
}

func (m *mockStudyEventService) GetByID(ctx context.Context, id studyevent.ID) (*studyevent.Event, error) {
	return m.getByIDFunc(ctx, id)
}

func (m *mockStudyEventService) ListRecentByUserAndCard(ctx context.Context, userID, cardID string, limit int) ([]*studyevent.Event, error) {
	return m.listRecentFunc(ctx, userID, cardID, limit)
}

func TestHandler_HealthCheck(t *testing.T) {
	h := NewHandler(nil, nil, nil, nil, nil, nil)
	resp, err := h.HealthCheck(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Status.Value)
}

func TestHandler_CreateUser(t *testing.T) {
	tests := []struct {
		name      string
		req       *oas.CreateUserRequest
		mockSetup func() UserService
		wantErr   bool
	}{
		{
			name: "success",
			req:  &oas.CreateUserRequest{Nickname: "test", ProviderID: "provider-1"},
			mockSetup: func() UserService {
				return &mockUserService{
					createFunc: func(ctx context.Context, nickname, providerID string) (*user.User, error) {
						u := user.New(nickname, providerID)
						return u, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "service error",
			req:  &oas.CreateUserRequest{Nickname: "test", ProviderID: "provider-1"},
			mockSetup: func() UserService {
				return &mockUserService{
					createFunc: func(ctx context.Context, nickname, providerID string) (*user.User, error) {
						return nil, errs.New(errs.ErrInternal, "db error")
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.mockSetup(), nil, nil, nil, nil, nil)
			_, err := h.CreateUser(context.Background(), tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestHandler_GetUser(t *testing.T) {
	tests := []struct {
		name       string
		params     oas.GetUserParams
		mockSetup  func() UserService
		wantErr    bool
		isNotFound bool
	}{
		{
			name:   "found",
			params: oas.GetUserParams{ID: "user-1"},
			mockSetup: func() UserService {
				return &mockUserService{
					getFunc: func(ctx context.Context, id user.ID) (*user.User, error) {
						return user.New("test", "provider-1"), nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "not found",
			params: oas.GetUserParams{ID: "user-1"},
			mockSetup: func() UserService {
				return &mockUserService{
					getFunc: func(ctx context.Context, id user.ID) (*user.User, error) {
						return nil, errs.New(errs.ErrNotFound, "not found")
					},
				}
			},
			wantErr:    false,
			isNotFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.mockSetup(), nil, nil, nil, nil, nil)
			res, err := h.GetUser(context.Background(), tt.params)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.isNotFound {
				_, ok := res.(*oas.GetUserNotFound)
				assert.True(t, ok, "expected NotFound response")
			}
		})
	}
}

func TestHandler_ListUsers(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() UserService
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func() UserService {
				return &mockUserService{
					listFunc: func(ctx context.Context) ([]*user.User, error) {
						return []*user.User{user.New("test", "provider-1")}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "service error",
			mockSetup: func() UserService {
				return &mockUserService{
					listFunc: func(ctx context.Context) ([]*user.User, error) {
						return nil, errs.New(errs.ErrInternal, "db error")
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.mockSetup(), nil, nil, nil, nil, nil)
			_, err := h.ListUsers(context.Background())
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestHandler_CreateCard(t *testing.T) {
	tests := []struct {
		name      string
		req       *oas.CreateCardRequest
		mockSetup func() CardService
		wantErr   bool
	}{
		{
			name: "success",
			req: &oas.CreateCardRequest{
				Question: oas.NewOptString("Q1"),
				Hint:     oas.NewOptString("hint"),
				Content:  "content",
				Creator:  "user-1",
			},
			mockSetup: func() CardService {
				return &mockCardService{
					createFunc: func(ctx context.Context, question, hint, content, creator string) (*card.Card, error) {
						return card.New(question, hint, content, creator), nil
					},
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(nil, tt.mockSetup(), nil, nil, nil, nil)
			_, err := h.CreateCard(context.Background(), tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestHandler_GetCard(t *testing.T) {
	tests := []struct {
		name       string
		params     oas.GetCardParams
		mockSetup  func() CardService
		wantErr    bool
		isNotFound bool
	}{
		{
			name:   "found",
			params: oas.GetCardParams{ID: "card-1"},
			mockSetup: func() CardService {
				return &mockCardService{
					getFunc: func(ctx context.Context, id card.ID) (*card.Card, error) {
						return card.New("Q1", "hint", "content", "user-1"), nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "not found",
			params: oas.GetCardParams{ID: "card-1"},
			mockSetup: func() CardService {
				return &mockCardService{
					getFunc: func(ctx context.Context, id card.ID) (*card.Card, error) {
						return nil, errs.New(errs.ErrNotFound, "not found")
					},
				}
			},
			wantErr:    false,
			isNotFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(nil, tt.mockSetup(), nil, nil, nil, nil)
			res, err := h.GetCard(context.Background(), tt.params)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.isNotFound {
				_, ok := res.(*oas.GetCardNotFound)
				assert.True(t, ok)
			}
		})
	}
}

func TestHandler_CreateStudyEvent(t *testing.T) {
	tests := []struct {
		name      string
		req       *oas.CreateStudyEventRequest
		mockSetup func() StudyEventService
		wantErr   bool
	}{
		{
			name: "success with quality",
			req: &oas.CreateStudyEventRequest{
				CollectionID: "collection-1",
				CardID:       "card-1",
				UserID:       "user-1",
				EventType:    "review",
				Quality:      oas.NewOptInt(4),
			},
			mockSetup: func() StudyEventService {
				return &mockStudyEventService{
					createFunc: func(ctx context.Context, collectionID, cardID, userID string, eventType studyevent.EventType, quality *int, responseTime *time.Duration, previousDeckID *string) (*studyevent.Event, error) {
						return studyevent.New(studyevent.ID("event-1"), collectionID, cardID, userID, eventType, quality, nil, nil), nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "success without quality",
			req: &oas.CreateStudyEventRequest{
				CollectionID: "collection-1",
				CardID:       "card-1",
				UserID:       "user-1",
				EventType:    "hint_shown",
			},
			mockSetup: func() StudyEventService {
				return &mockStudyEventService{
					createFunc: func(ctx context.Context, collectionID, cardID, userID string, eventType studyevent.EventType, quality *int, responseTime *time.Duration, previousDeckID *string) (*studyevent.Event, error) {
						return studyevent.New(studyevent.ID("event-1"), collectionID, cardID, userID, eventType, quality, nil, nil), nil
					},
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(nil, nil, nil, nil, tt.mockSetup(), nil)
			_, err := h.CreateStudyEvent(context.Background(), tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestHandler_GetNextCard_NoService(t *testing.T) {
	h := NewHandler(nil, nil, nil, nil, nil, nil)
	_, err := h.GetNextCard(context.Background(), oas.GetNextCardParams{UserID: "user-1"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "next card service not available")
}

func TestHandler_UpdateUser(t *testing.T) {
	tests := []struct {
		name       string
		req        *oas.UpdateUserRequest
		params     oas.UpdateUserParams
		mockSetup  func() UserService
		wantErr    bool
		isNotFound bool
	}{
		{
			name:   "success",
			req:    &oas.UpdateUserRequest{Nickname: "new-name", ProviderID: "provider-1"},
			params: oas.UpdateUserParams{ID: "user-1"},
			mockSetup: func() UserService {
				return &mockUserService{
					updateFunc: func(ctx context.Context, id user.ID, nickname, providerID string) (*user.User, error) {
						u := user.New(nickname, providerID)
						return u, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "not found",
			req:    &oas.UpdateUserRequest{Nickname: "new-name", ProviderID: "provider-1"},
			params: oas.UpdateUserParams{ID: "user-1"},
			mockSetup: func() UserService {
				return &mockUserService{
					updateFunc: func(ctx context.Context, id user.ID, nickname, providerID string) (*user.User, error) {
						return nil, errs.New(errs.ErrNotFound, "not found")
					},
				}
			},
			wantErr:    false,
			isNotFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.mockSetup(), nil, nil, nil, nil, nil)
			res, err := h.UpdateUser(context.Background(), tt.req, tt.params)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.isNotFound {
				_, ok := res.(*oas.UpdateUserNotFound)
				assert.True(t, ok)
			}
		})
	}
}

func TestHandler_DeleteUser(t *testing.T) {
	tests := []struct {
		name       string
		params     oas.DeleteUserParams
		mockSetup  func() UserService
		wantErr    bool
		isNotFound bool
	}{
		{
			name:   "success",
			params: oas.DeleteUserParams{ID: "user-1"},
			mockSetup: func() UserService {
				return &mockUserService{
					deleteFunc: func(ctx context.Context, id user.ID) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "not found",
			params: oas.DeleteUserParams{ID: "user-1"},
			mockSetup: func() UserService {
				return &mockUserService{
					deleteFunc: func(ctx context.Context, id user.ID) error {
						return errs.New(errs.ErrNotFound, "not found")
					},
				}
			},
			wantErr:    false,
			isNotFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.mockSetup(), nil, nil, nil, nil, nil)
			res, err := h.DeleteUser(context.Background(), tt.params)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.isNotFound {
				_, ok := res.(*oas.DeleteUserNotFound)
				assert.True(t, ok)
			}
		})
	}
}

func TestHandler_CreateDeck(t *testing.T) {
	tests := []struct {
		name      string
		req       *oas.CreateDeckRequest
		mockSetup func() DeckService
		wantErr   bool
	}{
		{
			name: "success",
			req:  &oas.CreateDeckRequest{Name: "Deck 1", Creator: "user-1"},
			mockSetup: func() DeckService {
				return &mockDeckService{
					createFunc: func(ctx context.Context, name, creator string) (*deck.Deck, error) {
						return deck.New(name, creator), nil
					},
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(nil, nil, tt.mockSetup(), nil, nil, nil)
			_, err := h.CreateDeck(context.Background(), tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestHandler_GetDeck(t *testing.T) {
	tests := []struct {
		name       string
		params     oas.GetDeckParams
		mockSetup  func() DeckService
		wantErr    bool
		isNotFound bool
	}{
		{
			name:   "found",
			params: oas.GetDeckParams{ID: "deck-1"},
			mockSetup: func() DeckService {
				return &mockDeckService{
					getFunc: func(ctx context.Context, id deck.ID) (*deck.Deck, error) {
						return deck.New("Deck 1", "user-1"), nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "not found",
			params: oas.GetDeckParams{ID: "deck-1"},
			mockSetup: func() DeckService {
				return &mockDeckService{
					getFunc: func(ctx context.Context, id deck.ID) (*deck.Deck, error) {
						return nil, errs.New(errs.ErrNotFound, "not found")
					},
				}
			},
			wantErr:    false,
			isNotFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(nil, nil, tt.mockSetup(), nil, nil, nil)
			res, err := h.GetDeck(context.Background(), tt.params)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.isNotFound {
				_, ok := res.(*oas.GetDeckNotFound)
				assert.True(t, ok)
			}
		})
	}
}

func TestHandler_CreateCollection(t *testing.T) {
	tests := []struct {
		name      string
		req       *oas.CreateCollectionRequest
		mockSetup func() CollectionService
		wantErr   bool
	}{
		{
			name: "success",
			req:  &oas.CreateCollectionRequest{Name: "Collection 1", Creator: "user-1"},
			mockSetup: func() CollectionService {
				return &mockCollectionService{
					createFunc: func(ctx context.Context, name, creator string) (*collection.Collection, error) {
						return collection.New(name, creator), nil
					},
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(nil, nil, nil, tt.mockSetup(), nil, nil)
			_, err := h.CreateCollection(context.Background(), tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestHandler_GetCollection(t *testing.T) {
	tests := []struct {
		name       string
		params     oas.GetCollectionParams
		mockSetup  func() CollectionService
		wantErr    bool
		isNotFound bool
	}{
		{
			name:   "found",
			params: oas.GetCollectionParams{ID: "collection-1"},
			mockSetup: func() CollectionService {
				return &mockCollectionService{
					getFunc: func(ctx context.Context, id collection.ID) (*collection.Collection, error) {
						return collection.New("Collection 1", "user-1"), nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "not found",
			params: oas.GetCollectionParams{ID: "collection-1"},
			mockSetup: func() CollectionService {
				return &mockCollectionService{
					getFunc: func(ctx context.Context, id collection.ID) (*collection.Collection, error) {
						return nil, errs.New(errs.ErrNotFound, "not found")
					},
				}
			},
			wantErr:    false,
			isNotFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(nil, nil, nil, tt.mockSetup(), nil, nil)
			res, err := h.GetCollection(context.Background(), tt.params)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.isNotFound {
				_, ok := res.(*oas.GetCollectionNotFound)
				assert.True(t, ok)
			}
		})
	}
}
