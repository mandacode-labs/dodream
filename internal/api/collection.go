package api

import (
	"context"

	"github.com/mandacode-labs/dodream/internal/collection"
	"github.com/mandacode-labs/dodream/internal/errs"
	"github.com/mandacode-labs/dodream/pkg/oas"
)

// CollectionService defines the interface for collection business logic.
type CollectionService interface {
	Create(ctx context.Context, name string, creator string) (*collection.Collection, error)
	GetByID(ctx context.Context, id collection.ID) (*collection.Collection, error)
	Update(ctx context.Context, id collection.ID, name string) (*collection.Collection, error)
	Delete(ctx context.Context, id collection.ID) error
	List(ctx context.Context) ([]*collection.Collection, error)
}

// CreateCollection implements createCollection operation.
func (h *Handler) CreateCollection(ctx context.Context, req *oas.CreateCollectionRequest) (*oas.Collection, error) {
	c, err := h.collectionService.Create(ctx, req.Name, req.Creator)
	if err != nil {
		return nil, err
	}
	return mapCollection(c), nil
}

// GetCollection implements getCollection operation.
func (h *Handler) GetCollection(ctx context.Context, params oas.GetCollectionParams) (oas.GetCollectionRes, error) {
	c, err := h.collectionService.GetByID(ctx, collection.ID(params.ID))
	if err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &oas.GetCollectionNotFound{}, nil
		}
		return nil, err
	}
	return mapCollection(c), nil
}

// ListCollections implements listCollections operation.
func (h *Handler) ListCollections(ctx context.Context) ([]oas.Collection, error) {
	collections, err := h.collectionService.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]oas.Collection, len(collections))
	for i, c := range collections {
		result[i] = *mapCollection(c)
	}
	return result, nil
}

// UpdateCollection implements updateCollection operation.
func (h *Handler) UpdateCollection(ctx context.Context, req *oas.UpdateCollectionRequest, params oas.UpdateCollectionParams) (oas.UpdateCollectionRes, error) {
	c, err := h.collectionService.Update(ctx, collection.ID(params.ID), req.Name)
	if err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &oas.UpdateCollectionNotFound{}, nil
		}
		return nil, err
	}
	return mapCollection(c), nil
}

// DeleteCollection implements deleteCollection operation.
func (h *Handler) DeleteCollection(ctx context.Context, params oas.DeleteCollectionParams) (oas.DeleteCollectionRes, error) {
	if err := h.collectionService.Delete(ctx, collection.ID(params.ID)); err != nil {
		if errs.Is(err, errs.ErrNotFound) {
			return &oas.DeleteCollectionNotFound{}, nil
		}
		return nil, err
	}
	return &oas.DeleteCollectionNoContent{}, nil
}
