package core

import "time"

// StudyEventID is a unique identifier for a StudyEvent.
type StudyEventID string

// String returns the string representation of the StudyEventID.
func (id StudyEventID) String() string {
	return string(id)
}

// StudyEventType represents the type of study event.
type StudyEventType string

const (
	// StudyEventTypeReview is emitted when the user reviews a card.
	StudyEventTypeReview StudyEventType = "review"
	// StudyEventTypeHintShown is emitted when the user reveals the hint.
	StudyEventTypeHintShown StudyEventType = "hint_shown"
	// StudyEventTypeCardOpened is emitted when the user opens a card.
	StudyEventTypeCardOpened StudyEventType = "card_opened"
	// StudyEventTypeSkipped is emitted when the user skips a card.
	StudyEventTypeSkipped StudyEventType = "skipped"
	// StudyEventTypeDeckChanged is emitted when the user moves a card between decks.
	StudyEventTypeDeckChanged StudyEventType = "deck_changed"
)

// StudyEvent represents a single study interaction event for event-sourced SRS.
type StudyEvent struct {
	id             StudyEventID
	collectionID   CollectionID
	cardID         CardID
	userID         UserID
	eventType      StudyEventType
	quality        *int           // 0~5 for review events
	responseTime   *time.Duration // time taken to answer
	previousDeckID *DeckID        // for deck_changed events
	createdAt      time.Time
}

// NewStudyEvent creates a new StudyEvent.
func NewStudyEvent(
	id StudyEventID,
	collectionID CollectionID,
	cardID CardID,
	userID UserID,
	eventType StudyEventType,
	quality *int,
	responseTime *time.Duration,
	previousDeckID *DeckID,
) *StudyEvent {
	return &StudyEvent{
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
func (se *StudyEvent) ID() StudyEventID {
	return se.id
}

// CollectionID returns the collection identifier.
func (se *StudyEvent) CollectionID() CollectionID {
	return se.collectionID
}

// CardID returns the card identifier.
func (se *StudyEvent) CardID() CardID {
	return se.cardID
}

// UserID returns the user identifier.
func (se *StudyEvent) UserID() UserID {
	return se.userID
}

// EventType returns the type of study event.
func (se *StudyEvent) EventType() StudyEventType {
	return se.eventType
}

// Quality returns the review quality (0~5), or nil for non-review events.
func (se *StudyEvent) Quality() *int {
	return se.quality
}

// ResponseTime returns the response duration, or nil if not recorded.
func (se *StudyEvent) ResponseTime() *time.Duration {
	return se.responseTime
}

// PreviousDeckID returns the previous deck ID for deck_changed events, or nil.
func (se *StudyEvent) PreviousDeckID() *DeckID {
	return se.previousDeckID
}

// CreatedAt returns the timestamp when the event was created.
func (se *StudyEvent) CreatedAt() time.Time {
	return se.createdAt
}
