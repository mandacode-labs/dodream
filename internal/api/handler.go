package api

import (
	"context"
	"fmt"
	"time"

	"github.com/mandacode-labs/dodream/internal/card"
	"github.com/mandacode-labs/dodream/internal/collection"
	"github.com/mandacode-labs/dodream/internal/deck"
	"github.com/mandacode-labs/dodream/internal/errs"
	"github.com/mandacode-labs/dodream/internal/redis"
	"github.com/mandacode-labs/dodream/internal/studyevent"
	"github.com/mandacode-labs/dodream/internal/user"
	gen "github.com/mandacode-labs/dodream/pkg/api/v1"
)

// Service interfaces consumed by the handler.
type UserService interface {
	Create(ctx context.Context, nickname string, providerID string) (*user.User, error)
	GetByID(ctx context.Context, id user.ID) (*user.User, error)
	Update(ctx context.Context, id user.ID, nickname string, providerID string) (*user.User, error)
	Delete(ctx context.Context, id user.ID) error
	List(ctx context.Context) ([]*user.User, error)
}

type CardService interface {
	Create(ctx context.Context, question string, hint string, content string, creator string) (*card.Card, error)
	GetByID(ctx context.Context, id card.ID) (*card.Card, error)
	Update(ctx context.Context, id card.ID, question string, hint string, content string) (*card.Card, error)
	Delete(ctx context.Context, id card.ID) error
	List(ctx context.Context) ([]*card.Card, error)
}

type DeckService interface {
	Create(ctx context.Context, name string, creator string) (*deck.Deck, error)
	GetByID(ctx context.Context, id deck.ID) (*deck.Deck, error)
	Update(ctx context.Context, id deck.ID, name string) (*deck.Deck, error)
	Delete(ctx context.Context, id deck.ID) error
	List(ctx context.Context) ([]*deck.Deck, error)
}

type CollectionService interface {
	Create(ctx context.Context, name string, creator string) (*collection.Collection, error)
	GetByID(ctx context.Context, id collection.ID) (*collection.Collection, error)
	Update(ctx context.Context, id collection.ID, name string) (*collection.Collection, error)
	Delete(ctx context.Context, id collection.ID) error
	List(ctx context.Context) ([]*collection.Collection, error)
}

type StudyEventService interface {
	CreateEvent(ctx context.Context, collectionID string, cardID string, userID string, eventType studyevent.EventType, quality *int, responseTime *time.Duration, previousDeckID *string) (*studyevent.Event, error)
	GetByID(ctx context.Context, id studyevent.ID) (*studyevent.Event, error)
	ListRecentByUserAndCard(ctx context.Context, userID string, cardID string, limit int) ([]*studyevent.Event, error)
}

// Handler implements the generated ogen Handler interface.
type Handler struct {
	userService       UserService
	cardService       CardService
	deckService       DeckService
	collectionService CollectionService
	studyEventService StudyEventService
	redisClient       *redis.Client
}

// NewHandler creates a new API handler.
func NewHandler(
	userService UserService,
	cardService CardService,
	deckService DeckService,
	collectionService CollectionService,
	studyEventService StudyEventService,
	redisClient *redis.Client,
) *Handler {
	return &Handler{
		userService:       userService,
		cardService:       cardService,
		deckService:       deckService,
		collectionService: collectionService,
		studyEventService: studyEventService,
		redisClient:       redisClient,
	}
}

// Ensure Handler implements the interface.
var _ gen.Handler = (*Handler)(nil)

// HealthCheck implements healthCheck operation.
func (h *Handler) HealthCheck(ctx context.Context) (*gen.HealthResponse, error) {
	return &gen.HealthResponse{Status: gen.NewOptString("ok")}, nil
}

// CreateUser implements createUser operation.
func (h *Handler) CreateUser(ctx context.Context, req *gen.CreateUserRequest) (*gen.User, error) {
	u, err := h.userService.Create(ctx, req.Nickname, req.ProviderID)
	if err != nil {
		return nil, err
	}
	return mapUser(u), nil
}

// GetUser implements getUser operation.
func (h *Handler) GetUser(ctx context.Context, params gen.GetUserParams) (gen.GetUserRes, error) {
	u, err := h.userService.GetByID(ctx, user.ID(params.ID))
	if err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &gen.GetUserNotFound{}, nil
		}
		return nil, err
	}
	return mapUser(u), nil
}

// ListUsers implements listUsers operation.
func (h *Handler) ListUsers(ctx context.Context) ([]gen.User, error) {
	users, err := h.userService.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]gen.User, len(users))
	for i, u := range users {
		result[i] = *mapUser(u)
	}
	return result, nil
}

// CreateCard implements createCard operation.
func (h *Handler) CreateCard(ctx context.Context, req *gen.CreateCardRequest) (*gen.Card, error) {
	c, err := h.cardService.Create(ctx, req.Question.Value, req.Hint.Value, req.Content, req.Creator)
	if err != nil {
		return nil, err
	}
	return mapCard(c), nil
}

// GetCard implements getCard operation.
func (h *Handler) GetCard(ctx context.Context, params gen.GetCardParams) (gen.GetCardRes, error) {
	c, err := h.cardService.GetByID(ctx, card.ID(params.ID))
	if err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &gen.GetCardNotFound{}, nil
		}
		return nil, err
	}
	return mapCard(c), nil
}

// ListCards implements listCards operation.
func (h *Handler) ListCards(ctx context.Context) ([]gen.Card, error) {
	cards, err := h.cardService.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]gen.Card, len(cards))
	for i, c := range cards {
		result[i] = *mapCard(c)
	}
	return result, nil
}

// CreateDeck implements createDeck operation.
func (h *Handler) CreateDeck(ctx context.Context, req *gen.CreateDeckRequest) (*gen.Deck, error) {
	d, err := h.deckService.Create(ctx, req.Name, req.Creator)
	if err != nil {
		return nil, err
	}
	return mapDeck(d), nil
}

// GetDeck implements getDeck operation.
func (h *Handler) GetDeck(ctx context.Context, params gen.GetDeckParams) (gen.GetDeckRes, error) {
	d, err := h.deckService.GetByID(ctx, deck.ID(params.ID))
	if err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &gen.GetDeckNotFound{}, nil
		}
		return nil, err
	}
	return mapDeck(d), nil
}

// ListDecks implements listDecks operation.
func (h *Handler) ListDecks(ctx context.Context) ([]gen.Deck, error) {
	decks, err := h.deckService.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]gen.Deck, len(decks))
	for i, d := range decks {
		result[i] = *mapDeck(d)
	}
	return result, nil
}

// CreateCollection implements createCollection operation.
func (h *Handler) CreateCollection(ctx context.Context, req *gen.CreateCollectionRequest) (*gen.Collection, error) {
	c, err := h.collectionService.Create(ctx, req.Name, req.Creator)
	if err != nil {
		return nil, err
	}
	return mapCollection(c), nil
}

// GetCollection implements getCollection operation.
func (h *Handler) GetCollection(ctx context.Context, params gen.GetCollectionParams) (gen.GetCollectionRes, error) {
	c, err := h.collectionService.GetByID(ctx, collection.ID(params.ID))
	if err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &gen.GetCollectionNotFound{}, nil
		}
		return nil, err
	}
	return mapCollection(c), nil
}

// ListCollections implements listCollections operation.
func (h *Handler) ListCollections(ctx context.Context) ([]gen.Collection, error) {
	collections, err := h.collectionService.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]gen.Collection, len(collections))
	for i, c := range collections {
		result[i] = *mapCollection(c)
	}
	return result, nil
}

// CreateStudyEvent implements createStudyEvent operation.
func (h *Handler) CreateStudyEvent(ctx context.Context, req *gen.CreateStudyEventRequest) (*gen.StudyEvent, error) {
	var quality *int
	if req.Quality.Set {
		q := req.Quality.Value
		quality = &q
	}

	var responseTime *time.Duration
	if req.ResponseTimeMs.Set {
		rt := time.Duration(req.ResponseTimeMs.Value) * time.Millisecond
		responseTime = &rt
	}

	var previousDeckID *string
	if req.PreviousDeckID.Set {
		pd := req.PreviousDeckID.Value
		previousDeckID = &pd
	}

	event, err := h.studyEventService.CreateEvent(
		ctx,
		req.CollectionID,
		req.CardID,
		req.UserID,
		studyevent.EventType(req.EventType),
		quality,
		responseTime,
		previousDeckID,
	)
	if err != nil {
		return nil, err
	}
	return mapStudyEvent(event), nil
}

// GetNextCard implements getNextCard operation.
func (h *Handler) GetNextCard(ctx context.Context, params gen.GetNextCardParams) (*gen.NextCardsResponse, error) {
	if h.redisClient == nil {
		return nil, fmt.Errorf("redis not available")
	}

	var limit int64 = 1
	if params.Limit.Set {
		limit = int64(params.Limit.Value)
	}

	collectionID := ""
	if params.CollectionID.Set {
		collectionID = params.CollectionID.Value
	}

	cards, err := h.redisClient.GetNextCards(ctx, params.UserID, collectionID, limit)
	if err != nil {
		return nil, err
	}

	results := make([]gen.CardRecommendation, len(cards))
	for i, c := range cards {
		cardID, _ := c.Member.(string)
		results[i] = gen.CardRecommendation{
			CardID:     gen.NewOptString(cardID),
			FinalScore: gen.NewOptFloat64(c.Score),
		}
	}

	return &gen.NextCardsResponse{Cards: results}, nil
}

// UpdateUser implements updateUser operation.
func (h *Handler) UpdateUser(ctx context.Context, req *gen.UpdateUserRequest, params gen.UpdateUserParams) (gen.UpdateUserRes, error) {
	u, err := h.userService.Update(ctx, user.ID(params.ID), req.Nickname, req.ProviderID)
	if err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &gen.UpdateUserNotFound{}, nil
		}
		return nil, err
	}
	return mapUser(u), nil
}

// DeleteUser implements deleteUser operation.
func (h *Handler) DeleteUser(ctx context.Context, params gen.DeleteUserParams) (gen.DeleteUserRes, error) {
	if err := h.userService.Delete(ctx, user.ID(params.ID)); err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &gen.DeleteUserNotFound{}, nil
		}
		return nil, err
	}
	return &gen.DeleteUserNoContent{}, nil
}

// UpdateCard implements updateCard operation.
func (h *Handler) UpdateCard(ctx context.Context, req *gen.UpdateCardRequest, params gen.UpdateCardParams) (gen.UpdateCardRes, error) {
	c, err := h.cardService.Update(ctx, card.ID(params.ID), req.Question.Value, req.Hint.Value, req.Content)
	if err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &gen.UpdateCardNotFound{}, nil
		}
		return nil, err
	}
	return mapCard(c), nil
}

// DeleteCard implements deleteCard operation.
func (h *Handler) DeleteCard(ctx context.Context, params gen.DeleteCardParams) (gen.DeleteCardRes, error) {
	if err := h.cardService.Delete(ctx, card.ID(params.ID)); err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &gen.DeleteCardNotFound{}, nil
		}
		return nil, err
	}
	return &gen.DeleteCardNoContent{}, nil
}

// UpdateDeck implements updateDeck operation.
func (h *Handler) UpdateDeck(ctx context.Context, req *gen.UpdateDeckRequest, params gen.UpdateDeckParams) (gen.UpdateDeckRes, error) {
	d, err := h.deckService.Update(ctx, deck.ID(params.ID), req.Name)
	if err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &gen.UpdateDeckNotFound{}, nil
		}
		return nil, err
	}
	return mapDeck(d), nil
}

// DeleteDeck implements deleteDeck operation.
func (h *Handler) DeleteDeck(ctx context.Context, params gen.DeleteDeckParams) (gen.DeleteDeckRes, error) {
	if err := h.deckService.Delete(ctx, deck.ID(params.ID)); err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &gen.DeleteDeckNotFound{}, nil
		}
		return nil, err
	}
	return &gen.DeleteDeckNoContent{}, nil
}

// UpdateCollection implements updateCollection operation.
func (h *Handler) UpdateCollection(ctx context.Context, req *gen.UpdateCollectionRequest, params gen.UpdateCollectionParams) (gen.UpdateCollectionRes, error) {
	c, err := h.collectionService.Update(ctx, collection.ID(params.ID), req.Name)
	if err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &gen.UpdateCollectionNotFound{}, nil
		}
		return nil, err
	}
	return mapCollection(c), nil
}

// DeleteCollection implements deleteCollection operation.
func (h *Handler) DeleteCollection(ctx context.Context, params gen.DeleteCollectionParams) (gen.DeleteCollectionRes, error) {
	if err := h.collectionService.Delete(ctx, collection.ID(params.ID)); err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &gen.DeleteCollectionNotFound{}, nil
		}
		return nil, err
	}
	return &gen.DeleteCollectionNoContent{}, nil
}
