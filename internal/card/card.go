// Package card provides the card domain.
package card

import (
	"time"

	"github.com/google/uuid"
)

// ID is a unique identifier for a Card.
type ID string

// String returns the string representation of the ID.
func (id ID) String() string {
	return string(id)
}

// Card represents a flashcard with question, hint, and content.
type Card struct {
	id        ID
	question  string
	hint      string
	content   string
	creatorID string // user ID as string to avoid import cycle
	createdAt time.Time
	updatedAt time.Time
}

// New creates a new Card with the given question, hint, content and creator.
func New(question string, hint string, content string, creatorID string) *Card {
	now := time.Now()
	return &Card{
		id:        ID(uuid.New().String()),
		question:  question,
		hint:      hint,
		content:   content,
		creatorID: creatorID,
		createdAt: now,
		updatedAt: now,
	}
}

// ID returns the unique identifier of the card.
func (c *Card) ID() ID {
	return c.id
}

// Question returns the card's question or prompt.
func (c *Card) Question() string {
	return c.question
}

// Hint returns the card's hint.
func (c *Card) Hint() string {
	return c.hint
}

// Content returns the card's main content or answer.
func (c *Card) Content() string {
	return c.content
}

// CreatorID returns the user ID who created this card.
func (c *Card) CreatorID() string {
	return c.creatorID
}

// CreatedAt returns the timestamp when the card was created.
func (c *Card) CreatedAt() time.Time {
	return c.createdAt
}

// UpdatedAt returns the timestamp when the card was last modified.
func (c *Card) UpdatedAt() time.Time {
	return c.updatedAt
}

// SetQuestion updates the card's question.
func (c *Card) SetQuestion(question string) {
	c.question = question
	c.updatedAt = time.Now()
}

// SetHint updates the card's hint.
func (c *Card) SetHint(hint string) {
	c.hint = hint
	c.updatedAt = time.Now()
}

// SetContent updates the card's content.
func (c *Card) SetContent(content string) {
	c.content = content
	c.updatedAt = time.Now()
}

// NewWithID creates a Card from database values.
func NewWithID(id ID, question string, hint string, content string, creatorID string, createdAt, updatedAt time.Time) *Card {
	return &Card{
		id:        id,
		question:  question,
		hint:      hint,
		content:   content,
		creatorID: creatorID,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}
