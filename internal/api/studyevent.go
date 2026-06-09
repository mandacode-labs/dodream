package api

import (
	"context"
	"time"

	"github.com/mandacode-labs/dodream/internal/studyevent"
	"github.com/mandacode-labs/dodream/pkg/oas"
)

// StudyEventService defines the interface for study event business logic.
type StudyEventService interface {
	CreateEvent(ctx context.Context, collectionID string, cardID string, userID string, eventType studyevent.EventType, quality *int, responseTime *time.Duration, previousDeckID *string) (*studyevent.Event, error)
	GetByID(ctx context.Context, id studyevent.ID) (*studyevent.Event, error)
	ListRecentByUserAndCard(ctx context.Context, userID string, cardID string, limit int) ([]*studyevent.Event, error)
}

// CreateStudyEvent implements createStudyEvent operation.
func (h *Handler) CreateStudyEvent(ctx context.Context, req *oas.CreateStudyEventRequest) (*oas.StudyEvent, error) {
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
