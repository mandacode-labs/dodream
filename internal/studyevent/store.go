package studyevent

import (
	"context"
	"time"

	"github.com/mandacode-labs/dodream/ent"
	entcard "github.com/mandacode-labs/dodream/ent/card"
	"github.com/mandacode-labs/dodream/ent/studyevent"
	"github.com/mandacode-labs/dodream/ent/user"
	"github.com/mandacode-labs/dodream/internal/errs"
)

// Store provides database operations for study events.
type Store struct {
	client *ent.Client
}

// NewStore creates a new Store with the given ent client.
func NewStore(client *ent.Client) *Store {
	return &Store{client: client}
}

// Compile-time check that Store implements Repository.
var _ Repository = (*Store)(nil)

func mapEntError(op string, err error) error {
	if ent.IsNotFound(err) {
		return errs.Wrap(errs.ErrNotFound, "not found", err)
	}
	if ent.IsConstraintError(err) {
		return errs.Wrap(errs.ErrConflict, "conflict", err)
	}
	return errs.Wrap(errs.ErrInternal, op, err)
}

// Create inserts a new study event into the database.
func (s *Store) Create(ctx context.Context, event *Event) (*Event, error) {
	builder := s.client.StudyEvent.Create().
		SetID(event.ID().String()).
		SetEventType(string(event.EventType())).
		SetCollectionID(event.CollectionID()).
		SetCardID(event.CardID()).
		SetUserID(event.UserID())

	if event.Quality() != nil {
		builder.SetQuality(*event.Quality())
	}
	if event.ResponseTime() != nil {
		builder.SetResponseTimeNs(event.ResponseTime().Nanoseconds())
	}
	if event.PreviousDeckID() != nil {
		builder.SetPreviousDeckID(*event.PreviousDeckID())
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return nil, mapEntError("create study event", err)
	}

	return s.toCore(created), nil
}

// GetByID retrieves a study event by ID.
func (s *Store) GetByID(ctx context.Context, id ID) (*Event, error) {
	event, err := s.client.StudyEvent.Query().
		Where(studyevent.ID(id.String())).
		Only(ctx)
	if err != nil {
		return nil, mapEntError("get study event by id", err)
	}
	return s.toCore(event), nil
}

// ListRecentByUserAndCard retrieves the most recent N study events for a user-card pair.
func (s *Store) ListRecentByUserAndCard(ctx context.Context, userID string, cardID string, limit int) ([]*Event, error) {
	events, err := s.client.StudyEvent.Query().
		Where(studyevent.And(
			studyevent.HasUserWith(user.ID(userID)),
			studyevent.HasCardWith(entcard.ID(cardID)),
		)).
		Order(ent.Desc(studyevent.FieldCreatedAt)).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, mapEntError("list recent study events", err)
	}

	result := make([]*Event, len(events))
	for i, event := range events {
		result[i] = s.toCore(event)
	}
	return result, nil
}

// toCore converts an ent StudyEvent to a domain Event.
func (s *Store) toCore(event *ent.StudyEvent) *Event {
	var quality *int
	if event.Quality != nil {
		q := *event.Quality
		quality = &q
	}

	var responseTime *time.Duration
	if event.ResponseTimeNs != nil {
		rt := time.Duration(*event.ResponseTimeNs)
		responseTime = &rt
	}

	var previousDeckID *string
	if event.Edges.PreviousDeck != nil {
		id := event.Edges.PreviousDeck.ID
		previousDeckID = &id
	}

	return New(
		ID(event.ID),
		event.Edges.Collection.ID,
		event.Edges.Card.ID,
		event.Edges.User.ID,
		EventType(event.EventType),
		quality,
		responseTime,
		previousDeckID,
	)
}
