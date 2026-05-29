package studyevent

import (
	"context"
	"time"

	"github.com/mandacode-labs/dodream/internal/errs"
	"github.com/mandacode-labs/dodream/internal/nats"
)

// Repository defines the storage operations for study events.
type Repository interface {
	Create(ctx context.Context, event *Event) (*Event, error)
	GetByID(ctx context.Context, id ID) (*Event, error)
	ListRecentByUserAndCard(ctx context.Context, userID string, cardID string, limit int) ([]*Event, error)
}

// Service provides study event business logic.
type Service struct {
	repo      Repository
	publisher *nats.Publisher
}

// NewService creates a new Service.
func NewService(repo Repository, publisher *nats.Publisher) *Service {
	return &Service{repo: repo, publisher: publisher}
}

// CreateEvent creates and publishes a study event.
func (s *Service) CreateEvent(ctx context.Context, collectionID string, cardID string, userID string, eventType EventType, quality *int, responseTime *time.Duration, previousDeckID *string) (*Event, error) {
	if quality != nil && (*quality < 0 || *quality > 5) {
		return nil, errs.New(errs.ErrInvalidInput, "quality must be between 0 and 5")
	}

	event := New(
		ID(""),
		collectionID,
		cardID,
		userID,
		eventType,
		quality,
		responseTime,
		previousDeckID,
	)

	saved, err := s.repo.Create(ctx, event)
	if err != nil {
		return nil, err
	}

	if s.publisher != nil {
		msg := &nats.StudyEventMessage{
			EventID:      saved.ID().String(),
			CollectionID: saved.CollectionID(),
			CardID:       saved.CardID(),
			UserID:       saved.UserID(),
			EventType:    string(saved.EventType()),
			CreatedAt:    saved.CreatedAt().Format("2006-01-02T15:04:05Z"),
		}
		if saved.Quality() != nil {
			q := *saved.Quality()
			msg.Quality = &q
		}
		if saved.ResponseTime() != nil {
			rt := saved.ResponseTime().Nanoseconds()
			msg.ResponseTimeNs = &rt
		}
		if saved.PreviousDeckID() != nil {
			pd := *saved.PreviousDeckID()
			msg.PreviousDeckID = &pd
		}
		_ = s.publisher.PublishStudyEvent(ctx, msg)
	}

	return saved, nil
}

// GetByID retrieves a study event by ID.
func (s *Service) GetByID(ctx context.Context, id ID) (*Event, error) {
	return s.repo.GetByID(ctx, id)
}

// ListRecentByUserAndCard retrieves recent events.
func (s *Service) ListRecentByUserAndCard(ctx context.Context, userID string, cardID string, limit int) ([]*Event, error) {
	return s.repo.ListRecentByUserAndCard(ctx, userID, cardID, limit)
}
