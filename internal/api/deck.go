package api

import (
	"context"

	"github.com/mandacode-labs/dodream/internal/deck"
	"github.com/mandacode-labs/dodream/internal/errs"
	"github.com/mandacode-labs/dodream/pkg/oas"
)

// DeckService defines the interface for deck business logic.
type DeckService interface {
	Create(ctx context.Context, name string, creator string) (*deck.Deck, error)
	GetByID(ctx context.Context, id deck.ID) (*deck.Deck, error)
	Update(ctx context.Context, id deck.ID, name string) (*deck.Deck, error)
	Delete(ctx context.Context, id deck.ID) error
	List(ctx context.Context) ([]*deck.Deck, error)
}

// CreateDeck implements createDeck operation.
func (h *Handler) CreateDeck(ctx context.Context, req *oas.CreateDeckRequest) (*oas.Deck, error) {
	d, err := h.deckService.Create(ctx, req.Name, req.Creator)
	if err != nil {
		return nil, err
	}
	return mapDeck(d), nil
}

// GetDeck implements getDeck operation.
func (h *Handler) GetDeck(ctx context.Context, params oas.GetDeckParams) (oas.GetDeckRes, error) {
	d, err := h.deckService.GetByID(ctx, deck.ID(params.ID))
	if err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &oas.GetDeckNotFound{}, nil
		}
		return nil, err
	}
	return mapDeck(d), nil
}

// ListDecks implements listDecks operation.
func (h *Handler) ListDecks(ctx context.Context) ([]oas.Deck, error) {
	decks, err := h.deckService.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]oas.Deck, len(decks))
	for i, d := range decks {
		result[i] = *mapDeck(d)
	}
	return result, nil
}

// UpdateDeck implements updateDeck operation.
func (h *Handler) UpdateDeck(ctx context.Context, req *oas.UpdateDeckRequest, params oas.UpdateDeckParams) (oas.UpdateDeckRes, error) {
	d, err := h.deckService.Update(ctx, deck.ID(params.ID), req.Name)
	if err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &oas.UpdateDeckNotFound{}, nil
		}
		return nil, err
	}
	return mapDeck(d), nil
}

// DeleteDeck implements deleteDeck operation.
func (h *Handler) DeleteDeck(ctx context.Context, params oas.DeleteDeckParams) (oas.DeleteDeckRes, error) {
	if err := h.deckService.Delete(ctx, deck.ID(params.ID)); err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &oas.DeleteDeckNotFound{}, nil
		}
		return nil, err
	}
	return &oas.DeleteDeckNoContent{}, nil
}
