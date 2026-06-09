package api

import (
	"context"
	"fmt"

	"github.com/mandacode-labs/dodream/pkg/oas"
)

// Handler implements the generated ogen Handler interface.
type Handler struct {
	userService       UserService
	cardService       CardService
	deckService       DeckService
	collectionService CollectionService
	studyEventService StudyEventService
	nextCardService   NextCardService
}

// NextCardService provides next card recommendations.
type NextCardService interface {
	GetNextCards(ctx context.Context, userID, collectionID string, limit int64) ([]oas.CardRecommendation, error)
}

// NewHandler creates a new API handler.
func NewHandler(
	userService UserService,
	cardService CardService,
	deckService DeckService,
	collectionService CollectionService,
	studyEventService StudyEventService,
	nextCardService NextCardService,
) *Handler {
	return &Handler{
		userService:       userService,
		cardService:       cardService,
		deckService:       deckService,
		collectionService: collectionService,
		studyEventService: studyEventService,
		nextCardService:   nextCardService,
	}
}

// Ensure Handler implements the interface.
var _ oas.Handler = (*Handler)(nil)

// HealthCheck implements healthCheck operation.
func (h *Handler) HealthCheck(ctx context.Context) (*oas.HealthResponse, error) {
	return &oas.HealthResponse{Status: oas.NewOptString("ok")}, nil
}

// GetNextCard implements getNextCard operation.
func (h *Handler) GetNextCard(ctx context.Context, params oas.GetNextCardParams) (*oas.NextCardsResponse, error) {
	if h.nextCardService == nil {
		return nil, fmt.Errorf("next card service not available")
	}

	var limit int64 = 1
	if params.Limit.Set {
		limit = int64(params.Limit.Value)
	}

	collectionID := ""
	if params.CollectionID.Set {
		collectionID = params.CollectionID.Value
	}

	cards, err := h.nextCardService.GetNextCards(ctx, params.UserID, collectionID, limit)
	if err != nil {
		return nil, err
	}

	return &oas.NextCardsResponse{Cards: cards}, nil
}
