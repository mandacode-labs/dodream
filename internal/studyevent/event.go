// Package studyevent provides the study event domain.
package studyevent

import "time"

// ID is a unique identifier for a StudyEvent.
type ID string

// String returns the string representation of the ID.
func (id ID) String() string {
	return string(id)
}

// EventType represents the type of study event.
type EventType string

const (
	// TypeReview is emitted when the user reviews a card.
	TypeReview EventType = "review"
	// TypeHintShown is emitted when the user reveals the hint.
	TypeHintShown EventType = "hint_shown"
	// TypeCardOpened is emitted when the user opens a card.
	TypeCardOpened EventType = "card_opened"
	// TypeSkipped is emitted when the user skips a card.
	TypeSkipped EventType = "skipped"
	// TypeDeckChanged is emitted when the user moves a card between decks.
	TypeDeckChanged EventType = "deck_changed"
)

// Event represents a single study interaction event for event-sourced SRS.
type Event struct {
	id             ID
	collectionID   string // as string to avoid import cycle
	cardID         string // as string to avoid import cycle
	userID         string // as string to avoid import cycle
	eventType      EventType
	quality        *int           // 0~5 for review events
	responseTime   *time.Duration // time taken to answer
	previousDeckID *string        // for deck_changed events
	createdAt      time.Time
}

// New creates a new StudyEvent.
func New(
	id ID,
	collectionID string,
	cardID string,
	userID string,
	eventType EventType,
	quality *int,
	responseTime *time.Duration,
	previousDeckID *string,
) *Event {
	return &Event{
		id:             id,
		collectionID:   collectionID,
		cardID:         cardID,
		userID:         userID,
		eventType:      eventType,
		quality:        quality,
		responseTime:   responseTime,
		previousDeckID: previousDeckID,
		createdAt:      time.Now(),
	}
}

// ID returns the unique identifier of the study event.
func (e *Event) ID() ID {
	return e.id
}

// CollectionID returns the collection identifier.
func (e *Event) CollectionID() string {
	return e.collectionID
}

// CardID returns the card identifier.
func (e *Event) CardID() string {
	return e.cardID
}

// UserID returns the user identifier.
func (e *Event) UserID() string {
	return e.userID
}

// EventType returns the type of study event.
func (e *Event) EventType() EventType {
	return e.eventType
}

// Quality returns the review quality (0~5), or nil for non-review events.
func (e *Event) Quality() *int {
	return e.quality
}

// ResponseTime returns the response duration, or nil if not recorded.
func (e *Event) ResponseTime() *time.Duration {
	return e.responseTime
}

// PreviousDeckID returns the previous deck ID for deck_changed events, or nil.
func (e *Event) PreviousDeckID() *string {
	return e.previousDeckID
}

// CreatedAt returns the timestamp when the event was created.
func (e *Event) CreatedAt() time.Time {
	return e.createdAt
}
