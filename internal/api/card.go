package api

import (
	"context"

	"github.com/mandacode-labs/dodream/internal/card"
	"github.com/mandacode-labs/dodream/internal/errs"
	"github.com/mandacode-labs/dodream/pkg/oas"
)

// CardService defines the interface for card business logic.
type CardService interface {
	Create(ctx context.Context, question string, hint string, content string, creator string) (*card.Card, error)
	GetByID(ctx context.Context, id card.ID) (*card.Card, error)
	Update(ctx context.Context, id card.ID, question string, hint string, content string) (*card.Card, error)
	Delete(ctx context.Context, id card.ID) error
	List(ctx context.Context) ([]*card.Card, error)
}

// CreateCard implements createCard operation.
func (h *Handler) CreateCard(ctx context.Context, req *oas.CreateCardRequest) (*oas.Card, error) {
	c, err := h.cardService.Create(ctx, req.Question.Value, req.Hint.Value, req.Content, req.Creator)
	if err != nil {
		return nil, err
	}
	return mapCard(c), nil
}

// GetCard implements getCard operation.
func (h *Handler) GetCard(ctx context.Context, params oas.GetCardParams) (oas.GetCardRes, error) {
	c, err := h.cardService.GetByID(ctx, card.ID(params.ID))
	if err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &oas.GetCardNotFound{}, nil
		}
		return nil, err
	}
	return mapCard(c), nil
}

// ListCards implements listCards operation.
func (h *Handler) ListCards(ctx context.Context) ([]oas.Card, error) {
	cards, err := h.cardService.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]oas.Card, len(cards))
	for i, c := range cards {
		result[i] = *mapCard(c)
	}
	return result, nil
}

// UpdateCard implements updateCard operation.
func (h *Handler) UpdateCard(ctx context.Context, req *oas.UpdateCardRequest, params oas.UpdateCardParams) (oas.UpdateCardRes, error) {
	c, err := h.cardService.Update(ctx, card.ID(params.ID), req.Question.Value, req.Hint.Value, req.Content)
	if err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &oas.UpdateCardNotFound{}, nil
		}
		return nil, err
	}
	return mapCard(c), nil
}

// DeleteCard implements deleteCard operation.
func (h *Handler) DeleteCard(ctx context.Context, params oas.DeleteCardParams) (oas.DeleteCardRes, error) {
	if err := h.cardService.Delete(ctx, card.ID(params.ID)); err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &oas.DeleteCardNotFound{}, nil
		}
		return nil, err
	}
	return &oas.DeleteCardNoContent{}, nil
}
